package assets

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/phoobynet/market-deck-server/clients"
	"github.com/phoobynet/market-deck-server/database"
	"github.com/schollz/closestmatch"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"sync"
)

var assetRepositoryOnce sync.Once

var assetRepository *Repository

type Repository struct {
	alpacaClient *alpaca.Client
	assets       cmap.ConcurrentMap[string, Asset]
	populated    bool
	db           *gorm.DB
	search       []string
	cm           *closestmatch.ClosestMatch
}

func GetRepository() *Repository {
	assetRepositoryOnce.Do(
		func() {
			assetRepository = &Repository{
				alpacaClient: clients.GetAlpacaClient(),
				assets:       cmap.New[Asset](),
				populated:    false,
				db:           database.GetDB(),
			}
		},
	)

	return assetRepository
}

func (a *Repository) Get(symbol string) *Asset {
	a.populate()

	if asset, ok := a.assets.Get(symbol); ok {
		return &asset
	}

	return nil
}

func (a *Repository) GetMulti(symbols []string) map[string]Asset {
	a.populate()

	var assets map[string]Asset

	a.db.Where("symbol IN ?", symbols).Find(&assets)

	return assets
}

func (a *Repository) GetAll() []Asset {
	a.populate()
	var assets []Asset

	a.db.Model(&Asset{}).Find(&assets)

	return assets
}

func (a *Repository) GetByClass(assetClass alpaca.AssetClass) []Asset {
	a.populate()

	var assets []Asset

	a.db.Model(&Asset{}).Find(&assets, "class = ?", assetClass)

	return assets
}

func (a *Repository) Search(searchPattern string) []Asset {
	a.populate()
	limit := 200
	assets := make([]Asset, 0)

	logrus.Printf("Searching for %s", searchPattern)

	results := a.cm.ClosestN(searchPattern, limit)
	logrus.Printf("Found %d results", len(results))

	possibleSymbol := strings.ToUpper(strings.Split(searchPattern, " ")[0])
	logrus.Printf("Possible symbol %s", possibleSymbol)
	var exactSymbolMatchAsset Asset
	var exactSymbolMatchFound bool

	for _, result := range results {
		symbol := strings.Split(result, " ")[0]

		if asset, ok := a.assets.Get(symbol); ok {
			if asset.Symbol == possibleSymbol {
				logrus.Printf("Exact symbol match found %s", possibleSymbol)
				exactSymbolMatchAsset = asset
				exactSymbolMatchFound = true
			} else {
				assets = append(assets, asset)
			}
		}
	}

	if exactSymbolMatchFound {
		assets = append([]Asset{exactSymbolMatchAsset}, assets...)
	}

	logrus.Printf("%v", assets)

	return assets
}

func (a *Repository) populate() {
	if a.populated {
		return
	}

	var count int64

	a.db.Model(&Asset{}).Count(&count)

	if count == 0 {
		alpacaAssets, err := a.alpacaClient.GetAssets(
			alpaca.GetAssetsRequest{
				Status:     string(alpaca.AssetActive),
				AssetClass: string(alpaca.USEquity),
			},
		)

		if err != nil {
			logrus.Panic(err)
		}

		var assets []Asset

		for _, alpacaAsset := range alpacaAssets {
			asset := FromAlpacaAsset(alpacaAsset)
			assets = append(assets, asset)
			a.assets.Set(alpacaAsset.Symbol, asset)
		}

		a.db.Create(assets)
	} else {
		var assets []Asset

		a.db.Model(&Asset{}).Find(&assets)

		for _, asset := range assets {
			a.assets.Set(asset.Symbol, asset)
		}
	}

	a.search = make([]string, a.assets.Count())

	a.assets.IterCb(
		func(key string, asset Asset) {
			a.search = append(a.search, fmt.Sprintf("%s %s", asset.Symbol, asset.Name))
		},
	)

	a.cm = closestmatch.New(a.search, []int{3})

	a.populated = true
}
