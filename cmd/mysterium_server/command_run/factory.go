package command_run

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/mysterium/node/communication"
	"github.com/mysterium/node/communication/nats_dialog"
	"github.com/mysterium/node/communication/nats_discovery"
	"github.com/mysterium/node/identity"
	"github.com/mysterium/node/ipify"
	"github.com/mysterium/node/nat"
	"github.com/mysterium/node/openvpn"
	openvpn_session "github.com/mysterium/node/openvpn/session"
	"github.com/mysterium/node/server"
	dto_discovery "github.com/mysterium/node/service_discovery/dto"
	"github.com/mysterium/node/session"
	"path/filepath"
)

func NewCommand(options CommandOptions) *CommandRun {
	return NewCommandWith(
		options,
		server.NewClient(),
		ipify.NewClient(),
		nat.NewService(),
	)
}

func NewCommandWith(
	options CommandOptions,
	mysteriumClient server.Client,
	ipifyClient ipify.Client,
	natService nat.NATService,
) *CommandRun {

	ks := keystore.NewKeyStore(options.DirectoryKeystore, keystore.StandardScryptN, keystore.StandardScryptP)
	identityHandler := NewNodeIdentityHandler(
		identity.NewIdentityManager(ks),
		mysteriumClient,
		options.DirectoryKeystore,
	)

	return &CommandRun{
		identitySelector: func() (identity.Identity, error) {
			return identityHandler.Select(options.NodeKey)
		},
		ipifyClient:     ipifyClient,
		mysteriumClient: mysteriumClient,
		natService:      natService,
		dialogWaiterFactory: func(identity identity.Identity) (communication.DialogWaiter, dto_discovery.Contact) {
			address := nats_discovery.NewAddressForIdentity(identity)
			return nats_dialog.NewDialogWaiter(address), address.GetContact()
		},
		sessionManagerFactory: func(vpnServerIp string) session.ManagerInterface {
			return openvpn_session.NewManager(openvpn.NewClientConfig(
				vpnServerIp,
				filepath.Join(options.DirectoryConfig, "ca.crt"),
				filepath.Join(options.DirectoryConfig, "client.crt"),
				filepath.Join(options.DirectoryConfig, "client.key"),
				filepath.Join(options.DirectoryConfig, "ta.key"),
			))
		},
		vpnServerFactory: func() *openvpn.Server {
			vpnServerConfig := openvpn.NewServerConfig(
				"10.8.0.0", "255.255.255.0",
				filepath.Join(options.DirectoryConfig, "ca.crt"),
				filepath.Join(options.DirectoryConfig, "server.crt"),
				filepath.Join(options.DirectoryConfig, "server.key"),
				filepath.Join(options.DirectoryConfig, "dh.pem"),
				filepath.Join(options.DirectoryConfig, "crl.pem"),
				filepath.Join(options.DirectoryConfig, "ta.key"),
			)
			return openvpn.NewServer(vpnServerConfig, options.DirectoryRuntime)
		},
	}
}
