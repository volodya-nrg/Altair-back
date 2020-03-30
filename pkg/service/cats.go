package service

import (
	"altair/api/response"
	"altair/pkg/helpers"
	"altair/server"
	"altair/storage"
	"reflect"
	"strconv"
	"sync"
)

func NewCatService() *CatService {
	return new(CatService)
}

type CatService struct{}

func (cs CatService) GetCats() ([]*storage.Cat, error) {
	cats := make([]*storage.Cat, 0)
	err := server.Db.Debug().Order("cat_id", false).Find(&cats).Error

	return cats, err
}
func (cs CatService) GetCatsFull() ([]*response.СatFull, error) {
	serviceProperties := NewPropertyService()
	catsFull := make([]*response.СatFull, 0)
	linkCatsProperties := make([]*storage.CatProperty, 0)
	catIds := make([]uint64, 0)

	cats, err := cs.GetCats()
	if err != nil {
		return catsFull, err
	}

	for _, cat := range cats {
		catFull := new(response.СatFull)
		catFull.Cat = cat
		catsFull = append(catsFull, catFull)
		catIds = append(catIds, cat.CatId)
	}

	propertiesFull, err := serviceProperties.GetPropertiesFullByCatIds(catIds)
	if err != nil {
		return catsFull, err
	}

	// надо подхватитть связи каталога со св-вами
	if err := server.Db.Debug().Table("cats_properties").Find(&linkCatsProperties).Error; err != nil {
		return catsFull, err
	}

	// имеем пулные каталоги, полные св-ва со значениями и их связи
	for _, catFull := range catsFull {
		for _, link := range linkCatsProperties {
			for _, prop := range propertiesFull {
				if link.CatId == catFull.CatId && link.PropertyId == prop.PropertyId {
					catFull.PropertiesFull = append(catFull.PropertiesFull, prop)
				}
			}
		}
	}

	return catsFull, nil
}
func (cs CatService) GetCatsFullAsTree(catsFull []*response.СatFull) *response.CatFullTree {
	treeFull := new(response.CatFullTree)

	for _, catFull := range catsFull {
		tmp := new(response.CatFullTree)
		tmp.СatFull = catFull

		if catFull.CatId > 0 {
			if catFull.ParentId == 0 {
				treeFull.Childes = append(treeFull.Childes, tmp)

			} else if catFull.ParentId > 0 {
				cs.createTreeFullWalk(treeFull, *tmp)
			}
		}
	}

	return treeFull
}
func (cs CatService) GetCatsAsTree(cats []*storage.Cat) *response.CatTree {
	tree := new(response.CatTree)

	for _, cat := range cats {
		if cat.CatId > 0 {
			tmp := new(response.CatTree)
			tmp.Cat = cat

			if cat.ParentId == 0 {
				tree.Childes = append(tree.Childes, tmp)

			} else if cat.ParentId > 0 {
				cs.createTreeWalk(tree, *tmp)
			}
		}
	}

	return tree
}
func (cs CatService) GetCatByID(catId uint64) (*storage.Cat, error) {
	pCat := new(storage.Cat)
	err := server.Db.Debug().First(pCat, catId).Error // проверяется в контроллере

	return pCat, err
}
func (cs CatService) GetCatFullByID(catId uint64) (*response.СatFull, error) {
	serviceProperties := NewPropertyService()
	pCatFull := new(response.СatFull)

	pCat, err := cs.GetCatByID(catId)
	if err != nil {
		return pCatFull, err
	}

	listPPropertiesFull, err := serviceProperties.GetPropertiesFullByCatId(catId)
	if err != nil {
		return pCatFull, err
	}

	pCatFull.Cat = pCat
	pCatFull.PropertiesFull = listPPropertiesFull

	return pCatFull, nil
}
func (cs CatService) Create(cat *storage.Cat) error {
	cat.Slug = helpers.TranslitRuToEn(cat.Name)

	if !server.Db.Debug().NewRecord(cat) {
		return errNotCreateNewCat
	}

	return server.Db.Debug().Create(cat).Error
}
func (cs CatService) Update(cat *storage.Cat) error {
	cat.Slug = helpers.TranslitRuToEn(cat.Name)

	return server.Db.Debug().Save(cat).Error
}
func (cs CatService) Delete(catId uint64) error {
	cat := storage.Cat{
		CatId: catId,
	}

	if err := server.Db.Debug().Delete(&cat).Error; err != nil {
		return err
	}
	if err := cs.deleteFromCatsPropertiesByCatId(catId); err != nil {
		return err
	}

	return nil
}
func (cs CatService) ReWriteCatsProperties(catId uint64, mPropertyId map[string]string, mPos map[string]string, mIsRequire map[string]string) ([]*storage.CatProperty, error) {
	list := make([]*storage.CatProperty, 0)
	tbl := server.Db.Debug().Table("cats_properties")

	if err := cs.deleteFromCatsPropertiesByCatId(catId); err != nil {
		return list, err
	}

	for k, sPropertyId := range mPropertyId {
		iPropertyId, err := strconv.ParseUint(sPropertyId, 10, 64)
		if err != nil {
			return list, err
		}

		catProperty := new(storage.CatProperty)
		catProperty.CatId = catId
		catProperty.PropertyId = iPropertyId

		if val, found := mPos[k]; found {
			if iPos, err := strconv.ParseUint(val, 10, 64); err == nil && iPos > 0 {
				catProperty.Pos = iPos
			}
		}

		if val, found := mIsRequire[k]; found {
			if val == "true" {
				catProperty.IsRequire = true
			}
		}

		if !tbl.NewRecord(catProperty) {
			return list, errNotCreateNewCatProperty
		}

		if err := tbl.Create(catProperty).Error; err != nil {
			return list, err
		}

		list = append(list, catProperty)
	}

	return list, nil
}
func (cs CatService) GetAncestorsNastedLoop(cats []storage.Cat, findCatId uint64) []storage.Cat {
	a := cs.ancestorsNastedLoopWalk(cats, findCatId, nil)

	// Reverse examples:

	// v1.
	// sort.Slice(a[:], func(i, j int) bool { return i > j })

	// v2.
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}

	// v3.
	//sort.Sort(ReverseCat(b))

	return a
}
func (cs CatService) GetDescendantsNastedLoop(catsTree *response.CatTree, findCatId uint64) *response.CatTree {
	result := new(response.CatTree)

	if findCatId == 0 {
		return catsTree
	}

	if !reflect.ValueOf(catsTree.Cat).IsNil() && catsTree.Cat.CatId == findCatId {
		return catsTree
	}

	for _, branch := range catsTree.Childes {
		if branch.Cat.CatId == findCatId {
			return branch

		} else if len(branch.Childes) > 0 {
			if res := cs.GetDescendantsNastedLoop(branch, findCatId); !reflect.ValueOf(res.Cat).IsNil() && res.Cat.CatId > 0 {
				return res
			}
		}
	}

	return result
}
func (cs CatService) GetDescendantsGoRutines(catsTree *response.CatTree, findCatId uint64) response.CatTree {
	var wg sync.WaitGroup
	out := response.CatTree{}

	for _, tree := range catsTree.Childes {
		wg.Add(1)
		go func(tmpTree response.CatTree) {
			defer wg.Done()
			out = cs.descendantsGoRutinesWalk(tmpTree, findCatId)
		}(*tree)
	}

	wg.Wait()

	return out
}
func (cs CatService) GetIdsFromCatsTree(catsTree *response.CatTree) []uint64 {
	result := make([]uint64, 0)
	uniq := make([]uint64, 0)

	if !reflect.ValueOf(catsTree.Cat).IsNil() {
		result = append(result, catsTree.CatId)
	}

	for _, v := range catsTree.Childes {
		if v.CatId > 0 {
			result = append(result, v.CatId)
		}
		if len(v.Childes) > 0 {
			result = append(result, cs.GetIdsFromCatsTree(v)...)
		}
	}

	// возьмем только уникальные значения
	for _, v1 := range result {
		has := false
		for _, v2 := range uniq {
			if v2 == v1 {
				has = true
			}
		}
		if !has {
			uniq = append(uniq, v1)
		}
	}

	return uniq
}

// private -------------------------------------------------------------------------------------------------------------
func (cs CatService) ancestorsNastedLoopWalk(cats []storage.Cat, findCatId uint64, receiver []storage.Cat) []storage.Cat {
	if receiver == nil {
		receiver = make([]storage.Cat, 0)
	}

	for _, cat := range cats {
		if cat.CatId == findCatId {
			receiver = append(receiver, cat)
			findCatId = cat.ParentId
			break
		}
	}

	if findCatId == 0 {
		return receiver
	}

	return cs.ancestorsNastedLoopWalk(cats, findCatId, receiver)
}
func (cs CatService) descendantsGoRutinesWalk(catTree response.CatTree, findCatId uint64) response.CatTree {
	result := response.CatTree{}

	if catTree.CatId == findCatId {
		return catTree
	}

	for _, tree := range catTree.Childes {
		if tree.CatId == findCatId {
			return *tree
		}

		if len(tree.Childes) > 0 {
			return cs.descendantsGoRutinesWalk(*tree, findCatId)
		}
	}

	return result
}
func (cs CatService) createTreeWalk(branches *response.CatTree, inputCat response.CatTree) {
	for _, branch := range branches.Childes {
		if branch.CatId == inputCat.ParentId {
			branch.Childes = append(branch.Childes, &inputCat)

		} else if len(branch.Childes) > 0 {
			cs.createTreeWalk(branch, inputCat)
		}
	}
}
func (cs CatService) createTreeFullWalk(branches *response.CatFullTree, inputCat response.CatFullTree) {
	for _, branch := range branches.Childes {
		if branch.CatId == inputCat.ParentId {
			branch.Childes = append(branch.Childes, &inputCat)

		} else if len(branch.Childes) > 0 {
			cs.createTreeFullWalk(branch, inputCat)
		}
	}
}
func (cs CatService) deleteFromCatsPropertiesByCatId(catId uint64) error {
	err := server.Db.Debug().
		Table("cats_properties").
		Where("cat_id = ?", catId).
		Delete(storage.CatProperty{}).Error

	return err
}

type ReverseCat []*storage.Cat

func (c ReverseCat) Len() int {
	return len(c)
}
func (c ReverseCat) Less(i, j int) bool {
	return i > j
}
func (c ReverseCat) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
