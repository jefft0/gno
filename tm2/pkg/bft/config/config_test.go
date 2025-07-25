package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// allows the default config to have a valid DB
	_ "github.com/gnolang/gno/tm2/pkg/db/goleveldb"
)

func TestConfig_LoadOrMakeConfigWithOptions(t *testing.T) {
	t.Parallel()

	t.Run("existing configuration", func(t *testing.T) {
		t.Parallel()

		// Provide an empty directory
		cfgDir := t.TempDir()
		cfgPath := filepath.Join(cfgDir, defaultConfigPath)

		// Create a default config
		cfg := DefaultConfig()
		cfg.SetRootDir(cfgDir)

		// Make an incremental changes
		cfg.Moniker = "custom moniker"

		// Make sure the cfg paths are initialized
		require.NoError(t, cfg.EnsureDirs())

		// Write the config
		require.NoError(t, WriteConfigFile(cfgPath, cfg))

		// Load the config
		loadedCfg, loadErr := LoadOrMakeConfigWithOptions(cfgDir)
		require.NoError(t, loadErr)

		assert.Equal(t, cfg, loadedCfg)
	})

	t.Run("no existing config", func(t *testing.T) {
		t.Parallel()

		// Provide an empty directory
		cfgDir := t.TempDir()
		cfgPath := filepath.Join(cfgDir, defaultConfigPath)

		cfg, err := LoadOrMakeConfigWithOptions(cfgDir)
		require.NoError(t, err)

		// Make sure the returned cfg is the default one
		expectedCfg := DefaultConfig()
		expectedCfg.SetRootDir(cfgDir)

		assert.Equal(t, expectedCfg, cfg)

		// Make sure the returned config was saved
		loadedCfg, loadErr := LoadConfigFile(cfgPath)
		require.NoError(t, loadErr)

		loadedCfg.SetRootDir(cfgDir)

		assert.Equal(t, cfg, loadedCfg)
	})

	t.Run("no existing config, with options", func(t *testing.T) {
		t.Parallel()

		moniker := "dummy moniker"

		// Provide an empty directory
		cfgDir := t.TempDir()
		cfgPath := filepath.Join(cfgDir, defaultConfigPath)

		cfg, err := LoadOrMakeConfigWithOptions(
			cfgDir,
			func(cfg *Config) {
				cfg.BaseConfig.Moniker = moniker
			},
		)
		require.NoError(t, err)

		// Make sure the returned config was saved
		loadedCfg, loadErr := LoadConfigFile(cfgPath)
		require.NoError(t, loadErr)

		loadedCfg.SetRootDir(cfgDir)

		assert.Equal(t, cfg, loadedCfg)
	})
}

func TestConfig_ValidateBaseConfig(t *testing.T) {
	t.Parallel()

	t.Run("valid default config", func(t *testing.T) {
		t.Parallel()

		c := DefaultConfig()

		assert.NoError(t, c.BaseConfig.ValidateBasic())
	})

	t.Run("invalid moniker", func(t *testing.T) {
		t.Parallel()

		c := DefaultConfig()
		c.Moniker = ""

		assert.ErrorIs(t, c.BaseConfig.ValidateBasic(), errInvalidMoniker)
	})

	t.Run("invalid DB backend", func(t *testing.T) {
		t.Parallel()

		c := DefaultConfig()
		c.DBBackend = "totally valid backend"

		assert.ErrorIs(t, c.BaseConfig.ValidateBasic(), errInvalidDBBackend)
	})

	t.Run("DB path not set", func(t *testing.T) {
		t.Parallel()

		c := DefaultConfig()
		c.DBPath = ""

		assert.ErrorIs(t, c.BaseConfig.ValidateBasic(), errInvalidDBPath)
	})

	t.Run("node key path not set", func(t *testing.T) {
		t.Parallel()

		c := DefaultConfig()
		c.NodeKey = ""

		assert.ErrorIs(t, c.BaseConfig.ValidateBasic(), errInvalidNodeKeyPath)
	})

	t.Run("invalid ABCI mechanism", func(t *testing.T) {
		t.Parallel()

		c := DefaultConfig()
		c.ABCI = "hopes and dreams"

		assert.ErrorIs(t, c.BaseConfig.ValidateBasic(), errInvalidABCIMechanism)
	})

	t.Run("invalid prof listen address", func(t *testing.T) {
		t.Parallel()

		c := DefaultConfig()
		c.ProfListenAddress = "beep.boop"

		assert.ErrorIs(t, c.BaseConfig.ValidateBasic(), errInvalidProfListenAddress)
	})
}

func TestConfig_DBDir(t *testing.T) {
	t.Parallel()

	t.Run("DB path is absolute", func(t *testing.T) {
		t.Parallel()

		c := DefaultConfig()
		c.RootDir = "/root"
		c.DBPath = "/abs/path"

		assert.Equal(t, c.DBPath, c.DBDir())
		assert.NotEqual(t, filepath.Join(c.RootDir, c.DBPath), c.DBDir())
	})

	t.Run("DB path is relative", func(t *testing.T) {
		t.Parallel()

		c := DefaultConfig()
		c.RootDir = "/root"
		c.DBPath = "relative/path"

		assert.Equal(t, filepath.Join(c.RootDir, c.DBPath), c.DBDir())
	})
}
