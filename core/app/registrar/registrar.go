package registrar

import (
	"context"
	"math"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"

	"github.com/aws/amazon-ssm-agent/agent/log"
	"github.com/aws/amazon-ssm-agent/common/identity"
	agentCtx "github.com/aws/amazon-ssm-agent/core/app/context"
)

const (
	minSleepSecondsBeforeRetry = 15
	jitterFactor               = 0.2
)

func getBackoffRetryJitterSleepDuration(retryCount int) time.Duration {
	// sleep for at least 15 seconds
	expBackoff := math.Max(minSleepSecondsBeforeRetry, math.Pow(2, float64(retryCount)))
	return time.Duration(int(expBackoff)+rand.Intn(int(math.Ceil(expBackoff*jitterFactor)))) * time.Second
}

type IRetryableRegistrar interface {
	Start() error
	Stop()
	GetRegistrationAttemptedChan() chan struct{}
}

type RetryableRegistrar struct {
	log                       log.T
	registrationAttemptedChan chan struct{}
	stopRegistrarChan         chan struct{}
	identityRegistrar         identity.Registrar
	timeAfterFunc             func(time.Duration) <-chan time.Time
	isRegistrarRunning        bool
	isRegistrarRunningLock    *sync.RWMutex
}

func NewRetryableRegistrar(agentCtx agentCtx.ICoreAgentContext) *RetryableRegistrar {
	log := agentCtx.Log().WithContext("[Registrar]")
	log.Debug("initializing registrar")
	// Cast to innerIdentityGetter interface that defined getInner
	innerGetter, ok := agentCtx.Identity().(identity.IInnerIdentityGetter)
	if !ok {
		log.Errorf("malformed identity")
		return nil
	}

	var identityRegistrar identity.Registrar
	if identityRegistrar, ok = innerGetter.GetInner().(identity.Registrar); !ok {
		log.Debug("identity does not leverage auto-registration")
		return nil
	}

	return &RetryableRegistrar{
		log:                       log,
		identityRegistrar:         identityRegistrar,
		registrationAttemptedChan: make(chan struct{}, 1),
		stopRegistrarChan:         make(chan struct{}),
		timeAfterFunc:             time.After,
		isRegistrarRunning:        false,
		isRegistrarRunningLock:    &sync.RWMutex{},
	}
}

func (r *RetryableRegistrar) Start() error {
	r.log.Info("Starting registrar module")
	r.setIsRegistrarRunning(true)
	go r.RegisterWithRetry()
	return nil
}

func (r *RetryableRegistrar) RegisterWithRetry() {
	defer func() {
		if err := recover(); err != nil {
			r.log.Errorf("registrar panic: %v", err)
			r.log.Errorf("Stacktrace:\n%s", debug.Stack())
			r.log.Flush()
			r.setIsRegistrarRunning(false)
			select {
			case <-r.registrationAttemptedChan:
				//channel open, write to channel to unblock and close
				r.registrationAttemptedChan <- struct{}{}
				close(r.registrationAttemptedChan)
			default:
			}
		}
	}()

	retryCount := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		errChan := make(chan error, 1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					r.log.Errorf("identity register panic: %v", err)
					r.log.Errorf("Stacktrace:\n%s", debug.Stack())
					r.log.Flush()
				}

				// Close errChan if still open
				select {
				case <-errChan:
				default:
					close(errChan)
				}
			}()

			errChan <- r.identityRegistrar.Register(ctx)
		}()
		select {
		case err := <-errChan:
			if retryCount == 0 {
				r.registrationAttemptedChan <- struct{}{}
				close(r.registrationAttemptedChan)
			}

			if err != nil {
				r.log.Errorf("failed to register identity: %v", err)
			} else {
				r.setIsRegistrarRunning(false)
				return
			}
		case <-r.stopRegistrarChan:
			cancel()
			r.log.Info("Stopping registrar")
			r.setIsRegistrarRunning(false)
			r.log.Flush()
			return
		}

		// Default sleep duration for non-aws errors
		sleepDuration := getBackoffRetryJitterSleepDuration(retryCount)
		// Max retry count is 16, which will sleep for about 18-22 hours
		if retryCount < 16 {
			retryCount++
		}

		r.log.Infof("sleeping for %v minutes before retrying registration", sleepDuration.Minutes())

		select {
		case <-r.stopRegistrarChan:
			cancel()
			r.log.Info("Stopping registrar")
			r.setIsRegistrarRunning(false)
			r.log.Flush()
			return
		case <-r.timeAfterFunc(sleepDuration):
		}
	}
}

func (r *RetryableRegistrar) setIsRegistrarRunning(isRegistrarRunning bool) {
	r.isRegistrarRunningLock.Lock()
	defer r.isRegistrarRunningLock.Unlock()

	r.isRegistrarRunning = isRegistrarRunning
}

func (r *RetryableRegistrar) getIsRegistrarRunning() bool {
	r.isRegistrarRunningLock.RLock()
	defer r.isRegistrarRunningLock.RUnlock()

	return r.isRegistrarRunning
}

// GetRegistrationAttemptedChan returns a channel that is written to and closed
// after registration is attempted or has succeeded
func (r *RetryableRegistrar) GetRegistrationAttemptedChan() chan struct{} {
	return r.registrationAttemptedChan
}

func (r *RetryableRegistrar) Stop() {
	if !r.getIsRegistrarRunning() {
		r.log.Info("Registrar is already stopped")
		r.log.Flush()
		return
	}

	r.log.Info("Sending signal to stop registrar")
	r.stopRegistrarChan <- struct{}{}
}
