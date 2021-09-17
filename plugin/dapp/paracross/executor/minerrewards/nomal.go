package minerrewards

import (
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/paracross/types"
)

type normal struct{}

func init() {
	register("normal", &normal{})
}

/ 
func (n *normal) GetConfigReward(cfg *types.Chain33Config, height int64) (int64, int64, int64) {
	coinReward := cfg.MGInt("mver.consensus.paracross.coinReward", height)
	fundReward := cfg.MGInt("mver.consensus.paracross.coinDevFund", height)
	coinBaseReward := cfg.MGInt("mver.consensus.paracross.coinBaseReward", height)

	if coinReward < 0 || fundReward < 0 || coinBaseReward < 0 {
		panic("para config consensus.paracross.coinReward should bigger than 0")
	}

	//decimalMode=false  1e8
	decimalMode := cfg.MIsEnable("mver.consensus.paracross.decimalMode", height)
	if !decimalMode {
		coinReward *= cfg.GetCoinPrecision()
		fundReward *= cfg.GetCoinPrecision()
		coinBaseReward *= cfg.GetCoinPrecision()
	}
	/ coinBaseReward ， coinBaseReward coinRewar 
	if coinBaseReward >= coinReward {
		coinBaseReward = coinReward / 10
	}
	return coinReward, fundReward, coinBaseReward
}

/ 
func (n *normal) RewardMiners(cfg *types.Chain33Config, coinReward int64, miners []string, height int64) ([]*pt.ParaMinerReward, int64) {
	/ 
	var change int64
	var rewards []*pt.ParaMinerReward
	/ 
	minerUnit := coinReward / int64(len(miners))
	if minerUnit > 0 {
		for _, m := range miners {
			r := &pt.ParaMinerReward{Addr: m, Amount: minerUnit}
			rewards = append(rewards, r)
		}

		/ 
		change = coinReward % minerUnit
	}
	return rewards, change
}
