package service

import (
	"altair/api/response"
	"altair/server"
	"altair/storage"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gorm.io/gorm"
	"strings"
)

// NewPropService - фабрика, создает объект свойства
func NewPropService() *PropService {
	return new(PropService)
}

// PropService - структура свойства
type PropService struct{}

// GetProps - получить все свойства
func (ps PropService) GetProps(order string) ([]*storage.Prop, error) {
	list := make([]*storage.Prop, 0)
	err := server.Db.Order(order).Find(&list).Error

	return list, err
}

// GetPropsWithKindName - получить свойства с именем (вид)
func (ps PropService) GetPropsWithKindName() ([]*response.PropWithKindName, error) {
	list := make([]*response.PropWithKindName, 0)
	query := `
		SELECT 	P.*, KP.name AS kind_prop_name
			FROM props P
				LEFT JOIN kind_props KP ON KP.kind_prop_id = P.kind_prop_id
			ORDER BY P.prop_id ASC`
	err := server.Db.Raw(query).Scan(&list).Error

	return list, err
}

// GetPropsFullByCatID - получить свойства (полные) относительно ID категории
func (ps PropService) GetPropsFullByCatID(catID uint64, withPropsOnlyFiltered bool) ([]*response.PropFull, error) {
	serviceValuesProp := NewValuesPropService()
	list := make([]*response.PropFull, 0)
	var sliceExpWhere = make([]string, 0)
	var where string

	if catID > 0 {
		sliceExpWhere = append(sliceExpWhere, fmt.Sprintf("%s%d", "CP.cat_id = ", catID))
	}
	if withPropsOnlyFiltered {
		sliceExpWhere = append(sliceExpWhere, "CP.is_can_as_filter = 1")
	}
	if len(sliceExpWhere) > 0 {
		where = "WHERE " + strings.Join(sliceExpWhere, " AND ")
	}

	query := `
		SELECT  P.*, 
				KP.name AS kind_prop_name, 
				CP.pos AS prop_pos,
				CP.is_require AS prop_is_require,
				CP.is_can_as_filter AS prop_is_can_as_filter,
				CP.comment AS prop_comment
			FROM props P
				LEFT JOIN kind_props KP ON KP.kind_prop_id = P.kind_prop_id
				LEFT JOIN cats_props CP ON CP.prop_id = P.prop_id
			` + where + `
			ORDER BY CP.pos ASC`

	if err := server.Db.Raw(query).Scan(&list).Error; err != nil {
		return list, err
	}
	if err := serviceValuesProp.PopulateWithValues(list); err != nil {
		return list, err
	}

	return list, nil
}

// GetPropsFullByCatIDs - получить свойства (полные) относительно нескольких ID категорий
func (ps PropService) GetPropsFullByCatIDs(catIDs []uint64) ([]*response.PropFull, error) {
	valuesPropService := NewValuesPropService()
	list := make([]*response.PropFull, 0)
	query := `
		SELECT  P.*,
				KP.name AS kind_prop_name, 
				CP.pos AS prop_pos,
				CP.is_require AS prop_is_require,
				CP.is_can_as_filter AS prop_is_can_as_filter,
				CP.comment AS prop_comment
			FROM props P
				LEFT JOIN kind_props KP ON KP.kind_prop_id = P.kind_prop_id
				LEFT JOIN cats_props CP ON CP.prop_id = P.prop_id
			WHERE CP.cat_id IN (?)
			ORDER BY CP.pos ASC`

	if err := server.Db.Raw(query, catIDs).Scan(&list).Error; err != nil {
		return list, err
	}
	if err := valuesPropService.PopulateWithValues(list); err != nil {
		return list, err
	}

	return list, nil
}

// GetPropByID - получить свойство относительно его ID
func (ps PropService) GetPropByID(propID uint64) (*storage.Prop, error) {
	prop := new(storage.Prop)
	err := server.Db.First(prop, propID).Error
	return prop, err
}

// GetPropFullByID - получить свойство (полное) относительно его ID
func (ps PropService) GetPropFullByID(propID uint64, valuesPropService *ValuesPropService) (*response.PropFull, error) {
	propFull := new(response.PropFull)

	// тут (в SQL) главно не проверить относительно категории
	query := `
		SELECT 	P.*, 
				KP.name AS kind_prop_name
			FROM props P
				LEFT JOIN kind_props KP ON KP.kind_prop_id = P.kind_prop_id
			WHERE P.prop_id = ?`
	err := server.Db.Raw(query, propID).Scan(propFull).Error // проверяется в контроллере.

	if propFull.Prop == nil {
		return propFull, gorm.ErrRecordNotFound
	}
	spew.Dump(propFull, err, errors.Is(err, gorm.ErrRecordNotFound))

	pValues, err := valuesPropService.GetValuesByPropID(propFull.PropID)
	if err != nil {
		return propFull, err
	}
	propFull.Values = pValues

	return propFull, err
}

// Create - создать свойство
func (ps PropService) Create(prop *storage.Prop, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Create(prop).Error

	return err
}

// Update - изменить свойство
func (ps PropService) Update(prop *storage.Prop, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db
	}

	err := tx.Save(prop).Error

	return err
}

// Delete - удалить свойство
func (ps PropService) Delete(propID uint64, tx *gorm.DB) error {
	valuesPropService := NewValuesPropService()

	if tx == nil {
		tx = server.Db
	}

	err := tx.Delete(storage.Prop{}, "prop_id = ?", propID).Error
	if err != nil {
		return err
	}

	err = valuesPropService.DeleteByPropID(propID, tx)
	if err != nil {
		return err
	}

	return err
}

// ReWriteValuesForProps - перезаписать значения для свойств
func (ps PropService) ReWriteValuesForProps(propID uint64, tx *gorm.DB, listNew []storage.ValueProp) ([]storage.ValueProp, error) {
	listOld := make([]storage.ValueProp, 0)
	listResult := make([]storage.ValueProp, 0)

	if tx == nil {
		tx = server.Db
	}

	// возьмем старый список
	if err := server.Db.Where("prop_id = ?", propID).Find(&listOld).Error; err != nil {
		return listResult, err
	}

	// возьмем id-шники которые надо удалить и при необходимости удалим
	removeID := make([]uint64, 0)
	for _, v1 := range listOld {
		isFound := false
		for _, v2 := range listNew {
			if v2.ValueID == v1.ValueID {
				isFound = true
			}
		}

		if !isFound {
			removeID = append(removeID, v1.ValueID)
		}
	}
	if len(removeID) > 0 {
		err := tx.Delete(storage.ValueProp{}, "value_id IN (?)", removeID).Error
		if err != nil {
			return listResult, err
		}
	}
	// \

	// обновим/добавим остальные эл-ты
	for _, v := range listNew {
		if v.ValueID == 0 {
			if err := tx.Create(&v).Error; err != nil {
				return listResult, err
			}

		} else if err := tx.Save(&v).Error; err != nil {
			return listResult, err
		}

		listResult = append(listResult, v)
	}

	return listResult, nil
}
