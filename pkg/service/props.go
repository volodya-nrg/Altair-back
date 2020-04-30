package service

import (
	"altair/api/response"
	"altair/server"
	"altair/storage"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func NewPropService() *PropService {
	return new(PropService)
}

type PropService struct{}

func (ps PropService) GetProps(order string) ([]*storage.Prop, error) {
	list := make([]*storage.Prop, 0)
	err := server.Db.Debug().Order(order).Find(&list).Error

	return list, err
}
func (ps PropService) GetPropsWithKindName() ([]*response.PropWithKindName, error) {
	list := make([]*response.PropWithKindName, 0)
	query := `
		SELECT 	P.*, KP.name AS kind_prop_name
			FROM props P
				LEFT JOIN kind_props KP ON KP.kind_prop_id = P.kind_prop_id
			ORDER BY P.prop_id ASC`
	err := server.Db.Debug().Raw(query).Scan(&list).Error

	return list, err
}
func (ps PropService) GetPropsFullByCatId(catId uint64, withPropsOnlyFiltered bool) ([]*response.PropFull, error) {
	serviceValuesProp := NewValuesPropService()
	list := make([]*response.PropFull, 0)
	var sliceExpWhere = make([]string, 0)
	var where string

	if catId > 0 {
		sliceExpWhere = append(sliceExpWhere, fmt.Sprintf("%s%d", "CP.cat_id = ", catId))
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

	if err := server.Db.Debug().Raw(query).Scan(&list).Error; err != nil {
		return list, err
	}

	if err := serviceValuesProp.PopulateWithValues(list); err != nil {
		return list, err
	}

	return list, nil
}
func (ps PropService) GetPropsFullByCatIds(catIds []uint64) ([]*response.PropFull, error) {
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

	if err := server.Db.Debug().Raw(query, catIds).Scan(&list).Error; err != nil {
		return list, err
	}

	if err := valuesPropService.PopulateWithValues(list); err != nil {
		return list, err
	}

	return list, nil
}
func (ps PropService) GetPropById(propId uint64) (*storage.Prop, error) {
	prop := new(storage.Prop)
	err := server.Db.Debug().First(prop, propId).Error
	return prop, err
}
func (ps PropService) GetPropFullById(propId uint64, valuesPropService *ValuesPropService) (*response.PropFull, error) {
	propFull := new(response.PropFull)
	query := `
		SELECT P.*, KP.name AS kind_prop_name 
			FROM props P
				LEFT JOIN kind_props KP ON KP.kind_prop_id = P.kind_prop_id
			WHERE P.prop_id = ?`
	err := server.Db.Debug().Raw(query, propId).Scan(&propFull).Error // проверяется в контроллере

	// добавим данные если есть куда
	if !gorm.IsRecordNotFoundError(err) {
		pValues, err := valuesPropService.GetValuesByPropId(propFull.PropId)
		if err != nil {
			return propFull, err
		}

		propFull.Values = pValues
	}

	return propFull, err
}
func (ps PropService) Create(prop *storage.Prop, tx *gorm.DB) error {
	if !server.Db.Debug().NewRecord(prop) {
		return errOnNewRecordNewProp
	}

	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Create(prop).Error

	return err
}
func (ps PropService) Update(prop *storage.Prop, tx *gorm.DB) error {
	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Save(prop).Error

	return err
}
func (ps PropService) Delete(propId uint64, tx *gorm.DB) error {
	valuesPropService := NewValuesPropService()

	if tx == nil {
		tx = server.Db.Debug()
	}

	err := tx.Delete(storage.Prop{}, "prop_id = ?", propId).Error
	if err != nil {
		return err
	}

	err = valuesPropService.DeleteByPropId(propId, tx)
	if err != nil {
		return err
	}

	return err
}
func (ps PropService) ReWriteValuesForProps(propId uint64, tx *gorm.DB,
	mId map[string]string, mTitle map[string]string, mPos map[string]string) ([]storage.ValueProp, error) {

	listOld := make([]storage.ValueProp, 0)
	listNew := make([]storage.ValueProp, 0)
	listResult := make([]storage.ValueProp, 0)

	if tx == nil {
		tx = server.Db.Debug()
	}

	// возьмем старый список
	if err := server.Db.Debug().Where("prop_id = ?", propId).Find(&listOld).Error; err != nil {
		return listResult, err
	}

	// создадим тут пришедший список
	for k, title := range mTitle {
		valueProp := storage.ValueProp{}
		valueProp.Title = strings.TrimSpace(title)
		valueProp.PropId = propId

		if val, found := mPos[k]; found {
			if iPos, err := strconv.ParseUint(val, 10, 64); err == nil && iPos > 0 {
				valueProp.Pos = iPos
			}
		}
		if val, found := mId[k]; found {
			if iId, err := strconv.ParseUint(val, 10, 64); err == nil && iId > 0 {
				valueProp.ValueId = iId
			}
		}

		listNew = append(listNew, valueProp)
	}
	// \

	// возьмем id-шники которые надо удалить и при необходимости удалим
	removeId := make([]uint64, 0)
	for _, v1 := range listOld {
		isFound := false
		for _, v2 := range listNew {
			if v2.ValueId == v1.ValueId {
				isFound = true
			}
		}

		if !isFound {
			removeId = append(removeId, v1.ValueId)
		}
	}
	if len(removeId) > 0 {
		err := tx.Delete(storage.ValueProp{}, "value_id IN (?)", removeId).Error
		if err != nil {
			return listResult, err
		}
	}
	// \

	// обновим/добавим остальные эл-ты
	for _, v := range listNew {
		if v.ValueId == 0 {
			if !server.Db.Debug().NewRecord(v) {
				return listResult, errOnNewRecordNewProp
			}

			if err := tx.Create(&v).Error; err != nil {
				return listResult, err
			}

		} else {
			if err := tx.Model(&v).Update(v).Error; err != nil {
				return listResult, err
			}
		}

		listResult = append(listResult, v)
	}

	return listResult, nil
}
