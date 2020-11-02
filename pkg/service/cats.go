package service

import (
	"altair/api/request"
	"altair/api/response"
	"altair/pkg/manager"
	"altair/server"
	"altair/storage"
	"gorm.io/gorm"
	"reflect"
)

// NewCatService - фабрика, создает объект категории
func NewCatService() *CatService {
	return new(CatService)
}

// CatService - структура категории
type CatService struct{}

// GetCats - получить все категории
func (cs CatService) GetCats(isDisabled int) ([]*storage.Cat, error) {
	cats := make([]*storage.Cat, 0)

	// порядок важен: parent_id asc, pos asc, cat_id asc
	stm := server.Db.Order("parent_id asc, pos asc, cat_id asc")

	if isDisabled > 0 {
		stm = stm.Where("is_disabled = 1")

	} else if isDisabled == 0 {
		stm = stm.Where("is_disabled = 0")
	}

	err := stm.Find(&cats).Error
	if err != nil {
		return cats, err
	}

	return cats, err
}

// GetCatsFullAsTree - получить полные данные о категорий в виде дерева
func (cs CatService) GetCatsFullAsTree(catsFull []*response.СatFull) *response.CatTreeFull {
	treeFull := new(response.CatTreeFull)

	for _, catFull := range catsFull {
		tmp := new(response.CatTreeFull)
		tmp.СatFull = catFull

		if catFull.CatID > 0 {
			if catFull.ParentID == 0 {
				treeFull.Childes = append(treeFull.Childes, tmp)

			} else if catFull.ParentID > 0 {
				cs.buildTreeFullWalk(treeFull, *tmp)
			}
		}
	}

	return treeFull
}

// GetCatsAsTree - получить категории в виде дерева
func (cs CatService) GetCatsAsTree(cats []*storage.Cat) *response.CatTree {
	tree := new(response.CatTree)
	tree.Cat = new(storage.Cat) // чтоб присутствовали и др. св-ва, по мимо childes

	for _, cat := range cats {
		if cat.CatID > 0 {
			tmp := new(response.CatTree)
			tmp.Cat = cat

			if cat.ParentID == 0 {
				tree.Childes = append(tree.Childes, tmp)

			} else if cat.ParentID > 0 {
				cs.buildTreeWalk(tree, *tmp)
			}
		}
	}

	return tree
}

// GetCatByID - получить данные об конкретной категории
func (cs CatService) GetCatByID(catID uint64, isDisabled int) (*storage.Cat, error) {
	cat := new(storage.Cat)
	stm := server.Db

	if isDisabled > 0 {
		stm = stm.Where("is_disabled = 1")

	} else if isDisabled == 0 {
		stm = stm.Where("is_disabled = 0")
	}

	err := stm.First(cat, catID).Error // проверяется в контроллере

	return cat, err
}

// GetCatFullByID - получить полные данные об конкретной категории
func (cs CatService) GetCatFullByID(catID uint64, withPropsOnlyFiltered bool, isDisabled int) (*response.СatFull, error) {
	serviceProps := NewPropService()
	catFull := new(response.СatFull)

	cat, err := cs.GetCatByID(catID, isDisabled)
	if err != nil {
		return catFull, err
	}

	propsFull, err := serviceProps.GetPropsFullByCatID(catID, withPropsOnlyFiltered)
	if err != nil {
		return catFull, err
	}

	catFull.Cat = cat
	catFull.PropsFull = propsFull

	return catFull, nil
}

// Create - создать категорию
func (cs CatService) Create(cat *storage.Cat, tx *gorm.DB) error {
	cat.Slug = manager.TranslitRuToEn(cat.Name)

	if tx == nil {
		tx = server.Db
	}

	err := tx.Create(cat).Error

	return err
}

// Update - изменить категорию
func (cs CatService) Update(cat *storage.Cat, tx *gorm.DB) error {
	cat.Slug = manager.TranslitRuToEn(cat.Name)

	if tx == nil {
		tx = server.Db
	}

	err := tx.Save(cat).Error

	return err
}

// Delete - удалить категорию
func (cs CatService) Delete(catID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}
	if err := tx.Where("cat_id = ?", catID).Delete(storage.Cat{}).Error; err != nil {
		return err
	}
	if err := cs.deleteFromCatsPropsByCatID(catID, tx); err != nil {
		return err
	}

	return nil
}

// ReWriteCatsProps - перезаписать свойства к конкретной категории
func (cs CatService) ReWriteCatsProps(catID uint64, tx *gorm.DB, propsAssignedForCat []request.PropAssignedForCat) ([]*storage.CatProp, error) {
	list := make([]*storage.CatProp, 0)

	if tx == nil {
		tx = server.Db
	}

	tbl := tx.Table("cats_props")

	if err := cs.deleteFromCatsPropsByCatID(catID, nil); err != nil {
		return list, err
	}

	for _, propAssignedForCat := range propsAssignedForCat {
		catProp := new(storage.CatProp)
		catProp.CatID = catID
		catProp.PropID = propAssignedForCat.PropID
		catProp.Pos = propAssignedForCat.Pos
		catProp.IsRequire = propAssignedForCat.IsRequire
		catProp.IsCanAsFilter = propAssignedForCat.IsCanAsFilter
		catProp.Comment = propAssignedForCat.Comment

		if err := tbl.Create(catProp).Error; err != nil {
			return list, err
		}

		list = append(list, catProp)
	}

	return list, nil
}

// GetAncestors - получить предков категории
func (cs CatService) GetAncestors(catsTree *response.CatTree, findCatID uint64) []storage.Cat { // предки
	list := make([]storage.Cat, 0)

	for _, branch := range catsTree.Childes {
		if branch.CatID == findCatID {
			list = append(list, *branch.Cat)
			return list
		}
		if len(branch.Childes) > 0 {
			res := cs.GetAncestors(branch, findCatID)

			if len(res) > 0 {
				list = append(list, *branch.Cat)
				list = append(list, res...)
				return list
			}
		}
	}

	return list
}

// GetDescendants - получить потомков категории
func (cs CatService) GetDescendants(catsTree *response.CatTree, findCatID uint64) *response.CatTree { // потомки
	result := new(response.CatTree)

	if findCatID == 0 {
		return catsTree
	}

	if !reflect.ValueOf(catsTree.Cat).IsNil() && catsTree.Cat.CatID == findCatID {
		return catsTree
	}

	for _, branch := range catsTree.Childes {
		if branch.Cat.CatID == findCatID {
			return branch

		} else if len(branch.Childes) > 0 {
			if res := cs.GetDescendants(branch, findCatID); !reflect.ValueOf(res.Cat).IsNil() && res.Cat.CatID > 0 {
				return res
			}
		}
	}

	return result
}

// GetIDsFromCatsTree - получить ID категорий из дерева категорий
func (cs CatService) GetIDsFromCatsTree(catsTree *response.CatTree) []uint64 {
	result := make([]uint64, 0)
	uniq := make([]uint64, 0)

	if !reflect.ValueOf(catsTree.Cat).IsNil() {
		result = append(result, catsTree.CatID)
	}

	for _, v := range catsTree.Childes {
		if v.CatID > 0 {
			result = append(result, v.CatID)
		}
		if len(v.Childes) > 0 {
			result = append(result, cs.GetIDsFromCatsTree(v)...)
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

// IsLeaf - является ли категория "листом" (конечной)
func (cs CatService) IsLeaf(catID uint64) (bool, error) {
	var count int64
	var result bool
	query := `
		SELECT COUNT(C.cat_id)
			FROM cats AS C
			WHERE 
        		C.cat_id = ?
            	AND (SELECT COUNT(*) FROM cats WHERE parent_id != C.cat_id AND cat_id = C.parent_id) > 0`

	if err := server.Db.Raw(query, catID).Count(&count).Error; err != nil {
		return result, err
	}

	result = count == 1

	return result, nil
}

// private -------------------------------------------------------------------------------------------------------------
func (cs CatService) buildTreeWalk(branches *response.CatTree, inputCat response.CatTree) {
	for _, branch := range branches.Childes {
		if branch.CatID == inputCat.ParentID {
			inputCat.Childes = make([]*response.CatTree, 0) // чтоб по умолчанию был []
			branch.Childes = append(branch.Childes, &inputCat)

		} else if len(branch.Childes) > 0 {
			cs.buildTreeWalk(branch, inputCat)
		}
	}
}
func (cs CatService) buildTreeFullWalk(branches *response.CatTreeFull, inputCat response.CatTreeFull) {
	for _, branch := range branches.Childes {
		if branch.CatID == inputCat.ParentID {
			inputCat.Childes = make([]*response.CatTreeFull, 0) // чтоб по умолчанию был []
			branch.Childes = append(branch.Childes, &inputCat)

		} else if len(branch.Childes) > 0 {
			cs.buildTreeFullWalk(branch, inputCat)
		}
	}
}
func (cs CatService) deleteFromCatsPropsByCatID(catID uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Table("cats_props").Delete(storage.CatProp{}, "cat_id = ?", catID).Error

	return err
}

// ReverseCat - сортировка списка (категорий)
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
