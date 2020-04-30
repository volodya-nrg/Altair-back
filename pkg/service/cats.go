package service

import (
	"altair/api/response"
	"altair/pkg/helpers"
	"altair/server"
	"altair/storage"
	"github.com/jinzhu/gorm"
	"reflect"
	"strconv"
	"strings"
)

func NewCatService() *CatService {
	return new(CatService)
}

type CatService struct{}

func (cs CatService) GetCats() ([]*storage.Cat, error) {
	cats := make([]*storage.Cat, 0)

	// порядок важен: parent_id asc, pos asc, cat_id asc
	err := server.Db.Debug().Order("parent_id asc, pos asc, cat_id asc").Find(&cats).Error
	if err != nil {
		return cats, err
	}

	return cats, err
}
func (cs CatService) GetCatsFullAsTree(catsFull []*response.СatFull) *response.CatTreeFull {
	treeFull := new(response.CatTreeFull)

	for _, catFull := range catsFull {
		tmp := new(response.CatTreeFull)
		tmp.СatFull = catFull

		if catFull.CatId > 0 {
			if catFull.ParentId == 0 {
				treeFull.Childes = append(treeFull.Childes, tmp)

			} else if catFull.ParentId > 0 {
				cs.buildTreeFullWalk(treeFull, *tmp)
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
				cs.buildTreeWalk(tree, *tmp)
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
func (cs CatService) GetCatFullByID(catId uint64, withPropsOnlyFiltered bool) (*response.СatFull, error) {
	serviceProps := NewPropService()
	catFull := new(response.СatFull)

	pCat, err := cs.GetCatByID(catId)
	if err != nil {
		return catFull, err
	}

	propsFull, err := serviceProps.GetPropsFullByCatId(catId, withPropsOnlyFiltered)
	if err != nil {
		return catFull, err
	}

	catFull.Cat = pCat
	catFull.PropsFull = propsFull

	return catFull, nil
}
func (cs CatService) Create(cat *storage.Cat, tx *gorm.DB) error {
	cat.Slug = helpers.TranslitRuToEn(cat.Name)

	if tx == nil {
		tx = server.Db.Debug()
	}
	if !server.Db.Debug().NewRecord(cat) {
		return errNotCreateNewCat
	}

	err := tx.Create(cat).Error

	return err
}
func (cs CatService) Update(cat *storage.Cat, tx *gorm.DB) error {
	cat.Slug = helpers.TranslitRuToEn(cat.Name)

	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Save(cat).Error

	return err
}
func (cs CatService) Delete(catId uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}
	if err := tx.Where("cat_id = ?", catId).Delete(storage.Cat{}).Error; err != nil {
		return err
	}
	if err := cs.deleteFromCatsPropsByCatId(catId, tx); err != nil {
		return err
	}

	return nil
}
func (cs CatService) ReWriteCatsProps(
	catId uint64,
	tx *gorm.DB,
	mPropId map[string]string,
	mPos map[string]string,
	mIsRequire map[string]string,
	mIsCanAsFilter map[string]string,
	mComment map[string]string) ([]*storage.CatProp, error) {

	list := make([]*storage.CatProp, 0)

	if tx == nil {
		tx = server.Db.Debug()
	}

	tbl := tx.Table("cats_props")

	if err := cs.deleteFromCatsPropsByCatId(catId, nil); err != nil {
		return list, err
	}

	for k, sPropId := range mPropId {
		iPropId, err := strconv.ParseUint(sPropId, 10, 64)
		if err != nil {
			return list, err
		}

		catProp := new(storage.CatProp)
		catProp.CatId = catId
		catProp.PropId = iPropId

		if val, found := mPos[k]; found {
			if iPos, err := strconv.ParseUint(val, 10, 64); err == nil && iPos > 0 {
				catProp.Pos = iPos
			}
		}

		if val, found := mIsRequire[k]; found {
			catProp.IsRequire = val == "true"
		}

		if val, found := mIsCanAsFilter[k]; found {
			catProp.IsCanAsFilter = val == "true"
		}

		if val, found := mComment[k]; found {
			catProp.Comment = strings.TrimSpace(val)
		}

		if !tbl.NewRecord(catProp) {
			return list, errNotCreateNewCatProp
		}

		if err := tbl.Create(catProp).Error; err != nil {
			return list, err
		}

		list = append(list, catProp)
	}

	return list, nil
}
func (cs CatService) GetAncestors(catsTree *response.CatTree, findCatId uint64) []storage.Cat { // предки
	list := make([]storage.Cat, 0)

	for _, branch := range catsTree.Childes {
		if branch.CatId == findCatId {
			list = append(list, *branch.Cat)
			return list
		}
		if len(branch.Childes) > 0 {
			res := cs.GetAncestors(branch, findCatId)

			if len(res) > 0 {
				list = append(list, *branch.Cat)
				list = append(list, res...)
				return list
			}
		}
	}

	return list
}
func (cs CatService) GetDescendants(catsTree *response.CatTree, findCatId uint64) *response.CatTree { // потомки
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
			if res := cs.GetDescendants(branch, findCatId); !reflect.ValueOf(res.Cat).IsNil() && res.Cat.CatId > 0 {
				return res
			}
		}
	}

	return result
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
func (cs CatService) buildTreeWalk(branches *response.CatTree, inputCat response.CatTree) {
	for _, branch := range branches.Childes {
		if branch.CatId == inputCat.ParentId {
			branch.Childes = append(branch.Childes, &inputCat)

		} else if len(branch.Childes) > 0 {
			cs.buildTreeWalk(branch, inputCat)
		}
	}
}
func (cs CatService) buildTreeFullWalk(branches *response.CatTreeFull, inputCat response.CatTreeFull) {
	for _, branch := range branches.Childes {
		if branch.CatId == inputCat.ParentId {
			branch.Childes = append(branch.Childes, &inputCat)

		} else if len(branch.Childes) > 0 {
			cs.buildTreeFullWalk(branch, inputCat)
		}
	}
}
func (cs CatService) deleteFromCatsPropsByCatId(catId uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Table("cats_props").Delete(storage.CatProp{}, "cat_id = ?", catId).Error

	return err
}

// ------------------
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
