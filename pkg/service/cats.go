package service

import (
	"altair/api/response"
	"altair/pkg/helpers"
	"altair/server"
	"altair/storage"
	"errors"
	"strconv"
	"sync"
)

var (
	errNotCreateNewCat         = errors.New("not create new cat")
	errNotCreateNewCatProperty = errors.New("not create new cat property")
)

func NewCatService() *CatService {
	return new(CatService)
}

type CatService struct{}

func (cs CatService) GetCats() ([]*storage.Cat, error) {
	cats := make([]*storage.Cat, 0)
	err := server.Db.Debug().Order("cat_id", false).Find(cats).Error

	return cats, err
}
func (cs CatService) GetCatsFull() ([]*response.СatFull, error) {
	//catsFull := make([]*response.СatFull, 0)
	//linkPropertiesWithCatsProperties := make([]*storage.LinkPropertiesWithCatsProperties, 0)
	//query := `
	//	SELECT P.*, CP.cat_id, CP.pos, CP.is_require
	//		FROM properties P
	//		LEFT JOIN cats_properties CP ON CP.property_id = P.property_id
	//		ORDER BY P.property_id ASC`
	//
	//if err := server.Db.Debug().Order("cat_id", false).Find(&catsFull).Error; err != nil {
	//	return catsFull, err
	//}
	//if err := server.Db.Debug().Raw(query).Scan(&linkPropertiesWithCatsProperties).Error; err != nil {
	//	return catsFull, err
	//}
	//
	//for _, cat := range catsFull {
	//	for _, link := range linkPropertiesWithCatsProperties {
	//		if cat.CatId == link.CatId {
	//			cat.PropertiesFull = append(cat.PropertiesFull, &link.PropertyFull)
	//		}
	//	}
	//}

	//logger.Info.Println(linkPropertiesWithCatsProperties)

	// заполняем св-вами, а так же их значениями
	// 1. берем список всего каталога + связи со св-вами
	// 2. берем список всех св-в + связи со значениями
	// 3. берем все значения
	// 4. все это дело соединить
	////////////////////////////////////////

	return []*response.СatFull{}, nil
}
func (cs CatService) GetCatsAsTree() (*response.CatTree, error) {
	cats := make([]*response.CatTree, 0)
	tree := new(response.CatTree)

	if err := server.Db.Debug().Table("cats").Order("parent_id, pos", false).Find(cats).Error; err != nil {
		return tree, err
	}

	for _, cat := range cats {
		if cat.ParentId == 0 {
			tree.Childes = append(tree.Childes, cat)

		} else if cat.ParentId > 0 {
			createTreeWalk(tree, *cat)
		}
	}

	return tree, nil
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
	if err := deleteFromCatsPropertiesByCatId(catId); err != nil {
		return err
	}

	return nil
}

func (cs CatService) ReWriteCatsProperties(catId uint64, mPropertyId map[string]string, mPos map[string]string, mIsRequire map[string]string) ([]*storage.CatProperty, error) {
	list := make([]*storage.CatProperty, 0)
	tbl := server.Db.Debug().Table("cats_properties")

	if err := deleteFromCatsPropertiesByCatId(catId); err != nil {
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
	a := ancestorsNastedLoopWalk(cats, findCatId, nil)

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
func (cs CatService) GetDescendantsNastedLoop(catsTree response.CatTree, findCatId uint64) response.CatTree {
	result := response.CatTree{}

	for _, branch := range catsTree.Childes {
		if branch.CatId == findCatId {
			result = *branch
			break

		} else if len(branch.Childes) > 0 {
			if res := cs.GetDescendantsNastedLoop(*branch, findCatId); res.CatId > 0 {
				result = res
				break
			}
		}
	}

	return result
}
func (cs CatService) GetDescendantsGoRutines(catsTree response.CatTree, findCatId uint64) response.CatTree {
	var wg sync.WaitGroup
	out := response.CatTree{}

	for _, tree := range catsTree.Childes {
		wg.Add(1)
		go func(tmpTree response.CatTree) {
			defer wg.Done()
			out = descendantsGoRutinesWalk(tmpTree, findCatId)
		}(*tree)
	}

	wg.Wait()

	return out
}
func (cs CatService) GetIdsFromCatsTree(catsTree *response.CatTree) []uint64 {
	result := make([]uint64, 0)

	for _, v := range catsTree.Childes {
		result = append(result, v.CatId)

		if len(v.Childes) > 0 {
			result = append(result, cs.GetIdsFromCatsTree(v)...)
		}
	}

	return result
}

// private -------------------------------------------------------------------------------------------------------------
func ancestorsNastedLoopWalk(cats []storage.Cat, findCatId uint64, receiver []storage.Cat) []storage.Cat {
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

	return ancestorsNastedLoopWalk(cats, findCatId, receiver)
}
func descendantsGoRutinesWalk(catTree response.CatTree, findCatId uint64) response.CatTree {
	result := response.CatTree{}

	if catTree.CatId == findCatId {
		return catTree
	}

	for _, tree := range catTree.Childes {
		if tree.CatId == findCatId {
			return *tree
		}

		if len(tree.Childes) > 0 {
			return descendantsGoRutinesWalk(*tree, findCatId)
		}
	}

	return result
}
func createTreeWalk(branches *response.CatTree, inputCat response.CatTree) {
	for _, branch := range branches.Childes {
		if branch.CatId == inputCat.ParentId {
			branch.Childes = append(branch.Childes, &inputCat)

		} else if len(branch.Childes) > 0 {
			createTreeWalk(branch, inputCat)
		}
	}
}
func deleteFromCatsPropertiesByCatId(catId uint64) error {
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
