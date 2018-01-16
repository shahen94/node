package command_run

import (
	"errors"
	"github.com/mysterium/node/identity"
	"github.com/mysterium/node/server"
)

//LoadIdentity selects and unlocks lastUsed identity or creates and unlocks new one if keyOption is not present
func LoadIdentity(identityHandler *identityHandler, keyOption, passphrase string) (id identity.Identity, err error) {
	if len(keyOption) > 0 {
		return identityHandler.UseExisting(keyOption, passphrase)
	}

	if id, err = identityHandler.UseLast(passphrase); err == nil {
		return id, err
	}

	return identityHandler.UseNew(passphrase)
}

type identityHandler struct {
	manager       identity.IdentityManagerInterface
	identityApi   server.Client
	cache         identity.IdentityCacheInterface
	signerFactory identity.SignerFactory
}

//NewNodeIdentityHandler creates new identity handler used by node
func NewNodeIdentityHandler(
	manager identity.IdentityManagerInterface,
	identityApi server.Client,
	cache identity.IdentityCacheInterface,
	signerFactory identity.SignerFactory,
) *identityHandler {
	return &identityHandler{
		manager:       manager,
		identityApi:   identityApi,
		cache:         cache,
		signerFactory: signerFactory,
	}
}

func (ih *identityHandler) UseExisting(address, passphrase string) (id identity.Identity, err error) {
	id, err = ih.manager.GetIdentity(address)
	if err != nil {
		return
	}
	err = ih.manager.Unlock(address, passphrase)
	if err != nil {
		return
	}

	err = ih.cache.StoreIdentity(id)
	return
}

func (ih *identityHandler) UseLast(passphrase string) (identity identity.Identity, err error) {
	identity, err = ih.cache.GetIdentity()
	if err != nil || !ih.manager.HasIdentity(identity.Address) {
		return identity, errors.New("identity not found in cache")
	}

	err = ih.manager.Unlock(identity.Address, passphrase)
	if err != nil {
		return identity, err
	}

	return identity, nil
}

func (ih *identityHandler) UseNew(passphrase string) (id identity.Identity, err error) {
	// if all fails, create a new one
	id, err = ih.manager.CreateNewIdentity(passphrase)
	if err != nil {
		return
	}

	err = ih.manager.Unlock(id.Address, passphrase)
	if err != nil {
		return
	}

	if err = ih.identityApi.RegisterIdentity(id, ih.signerFactory(id)); err != nil {
		return
	}

	err = ih.cache.StoreIdentity(id)
	return
}
