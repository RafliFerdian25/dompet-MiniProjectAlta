package subCategoryService

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/subCategoryRepository"
	"errors"
)

type SubCategoryService interface {
	GetSubCategoryByUser(userId uint) ([]dto.SubCategoryDTO, error)
	CreateSubCategory(subCategory dto.SubCategoryDTO) error
	DeleteSubCategory(id uint, userId uint) error
	UpdateSubCategory(subCategory dto.SubCategoryDTO) error
}

type subCategoryService struct {
	subCategoryRepository subCategoryRepository.SubCategoryRepository
}

// GetSubCategoryByUser implements SubCategoryService
func (sc *subCategoryService) GetSubCategoryByUser(userId uint) ([]dto.SubCategoryDTO, error) {
	subCategoryUsers, err := sc.subCategoryRepository.GetSubCategoryByUser(userId)
	if err != nil {
		return nil, err
	}
	return subCategoryUsers, nil
}

// CreateSubCategory implements SubCategoryService
func (sc *subCategoryService) CreateSubCategory(subCategory dto.SubCategoryDTO) error {
	err := sc.subCategoryRepository.CreateSubCategory(subCategory)
	if err != nil {
		return err
	}
	return nil
}

// DeleteSubCategory implements SubCategoryService
func (sc *subCategoryService) DeleteSubCategory(id uint, userId uint) error {
	subCategory, err := sc.subCategoryRepository.GetSubCategoryById(id)
	if err != nil {
		return err
	}
	// check if user id in the subcategory is the same as the user id in the token
	if *subCategory.UserID != userId {
		return errors.New(constantError.ErrorNotAuthorized)
	}
	err = sc.subCategoryRepository.DeleteSubCategory(id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSubCategory implements SubCategoryService
func (sc *subCategoryService) UpdateSubCategory(subCategory dto.SubCategoryDTO) error {
	dataSubCategory, err := sc.subCategoryRepository.GetSubCategoryById(subCategory.ID)
	if err != nil {
		return err
	}
	// check if user id in the subcategory is the same as the user id in the token
	if subCategory.UserID != dataSubCategory.UserID {
		return errors.New(constantError.ErrorNotAuthorized)
	}

	err = sc.subCategoryRepository.UpdateSubCategory(subCategory)
	if err != nil {
		return err
	}
	return nil
}

func NewSubCategoryService(subCategoryRepository subCategoryRepository.SubCategoryRepository) SubCategoryService {
	return &subCategoryService{
		subCategoryRepository: subCategoryRepository,
	}
}
