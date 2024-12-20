package solana

import (
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
)

type PrimorialEntry struct {
	Pubkey   string `pulumi:"pubkey"`
	Lamports string `pulumi:"lamports"`
}

// GenesisFlags represents the configuration flags for the Solana genesis setup.
type GenesisFlags struct {
	LedgerPath                 string    `pulumi:"ledgerPath"`
	IdentityPubkey             string    `pulumi:"identityPubkey"`
	VotePubkey                 string    `pulumi:"votePubkey"`
	StakePubkey                string    `pulumi:"stakePubkey"`
	ExtraFlags                 *[]string `pulumi:"extraFlags,optional"`
	FaucetPubkey               string    `pulumi:"faucetPubkey"`
	FaucetLamports             *string   `pulumi:"faucetLamports,optional"`
	TargetLamportsPerSignature *string   `pulumi:"targetLamportsPerSignature,optional"`
	Inflation                  *string   `pulumi:"inflation,optional"`
	LamportsPerByteYear        *string   `pulumi:"lamportsPerByteYear,optional"`
	SlotPerEpoch               *string   `pulumi:"slotPerEpoch,optional"`
	ClusterType                *string   `pulumi:"clusterType,optional"`
}

type CreateCommand struct {
	Genesis
}

func (cmd *CreateCommand) Check() error {
	return nil
}

func (cmd *CreateCommand) Env() *runner.EnvBuilder {
	b := runner.NewEnvBuilder()

	b.SetMap(map[string]string{
		"LEDGER_PATH":                   cmd.Flags.LedgerPath,
		"IDENTITY_PUBKEY":               cmd.Flags.IdentityPubkey,
		"VOTE_PUBKEY":                   cmd.Flags.VotePubkey,
		"STAKE_PUBKEY":                  cmd.Flags.StakePubkey,
		"FAUCET_PUBKEY":                 cmd.Flags.FaucetPubkey,
		"FAUCET_LAMPORTS":               "1000",
		"TARGET_LAMPORTS_PER_SIGNATURE": "0",
		"INFLATION":                     "none",
		"LAMPORTS_PER_BYTE_YEAR":        "1",
		"SLOT_PER_EPOCH":                "150",
		"CLUSTER_TYPE":                  "development",
	})

	b.SetArrayP("GENESIS_EXTRA_FLAGS", cmd.Flags.ExtraFlags)
	b.SetP("FAUCET_LAMPORTS", cmd.Flags.FaucetLamports)
	b.SetP("TARGET_LAMPORTS_PER_SIGNATURE", cmd.Flags.TargetLamportsPerSignature)
	b.SetP("INFLATION", cmd.Flags.Inflation)
	b.SetP("LAMPORTS_PER_BYTE_YEAR", cmd.Flags.LamportsPerByteYear)
	b.SetP("SLOT_PER_EPOCH", cmd.Flags.SlotPerEpoch)
	b.SetP("CLUSTER_TYPE", cmd.Flags.ClusterType)

	var primordialPubkeys, primordialLamports string
	if cmd.Primordial != nil {
		var pubkeys, lamports []string
		for _, entry := range cmd.Primordial {
			pubkeys = append(pubkeys, entry.Pubkey)
			lamports = append(lamports, entry.Lamports)
		}
		primordialPubkeys = strings.Join(pubkeys, ",")
		primordialLamports = strings.Join(lamports, ",")
	}

	b.Set("PRIMORDIAL_PUBKEYS", primordialPubkeys)
	b.Set("PRIMORDIAL_LAMPORTS", primordialLamports)

	b.SetP("PACKAGE_VERSION", cmd.Version)

	return b
}

func (cmd *CreateCommand) AddToPayload(p *runner.Payload) error {
	genesisScript, err := assets.Open(assetsGenesisScript)

	if err != nil {
		return err
	}

	p.AddReader("steps.sh", genesisScript)

	return nil
}

type Genesis struct {
	Flags      GenesisFlags     `pulumi:"flags"`
	Primordial []PrimorialEntry `pulumi:"primordial"`
	Version    *string          `pulumi:"version,optional"`
}

func (g *Genesis) Create() runner.Command {
	return &CreateCommand{
		Genesis: *g,
	}
}
