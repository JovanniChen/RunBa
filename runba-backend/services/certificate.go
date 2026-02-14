package services

import (
	"errors"
	"strings"

	"steam-backend/config"
	"steam-backend/models"

	"gorm.io/gorm"
)

// CertificateService 证书服务
// 处理证书相关的业务逻辑
// 说明：数据库表结构由外部维护，不依赖GORM自动迁移
//
// Table: certificates
// Model: models.Certificate
type CertificateService struct {
	db *gorm.DB
}

// NewCertificateService 创建新的证书服务实例
func NewCertificateService() *CertificateService {
	return &CertificateService{
		db: config.GetDB(),
	}
}

// CreateCertificate 创建证书
func (s *CertificateService) CreateCertificate(req *models.CreateCertificateRequest) (*models.Certificate, error) {
	if strings.TrimSpace(req.Number) == "" {
		return nil, errors.New("证书编号不能为空")
	}

	certificate := &models.Certificate{
		Number:           strings.TrimSpace(req.Number),
		OverallLength:    req.OverallLength,
		BladeLength:      req.BladeLength,
		HandleLength:     req.HandleLength,
		BladeMaterial:    req.BladeMaterial,
		BladeThickness:   req.BladeThickness,
		KissanLength:     req.KissanLength,
		Hamon:            req.Hamon,
		Mekugi:           req.Mekugi,
		Sakihaba:         req.Sakihaba,
		Motohaba:         req.Motohaba,
		SayaMaterial:     req.SayaMaterial,
		TsubaMaterial:    req.TsubaMaterial,
		HabakiMaterial:   req.HabakiMaterial,
		ItoSageoMaterial: req.ItoSageoMaterial,
		ForgeID:          req.ForgeID,
		SwordsmithID:     req.SwordsmithID,
		DateCompleted:    req.DateCompleted,
		No:               req.No,
	}

	if err := s.db.Create(certificate).Error; err != nil {
		return nil, err
	}

	return certificate, nil
}

// GetCertificates 获取证书列表
func (s *CertificateService) GetCertificates(params *models.CertificateListRequest) (*models.CertificateListResponse, error) {
	var certificates []models.Certificate
	var total int64

	query := s.db.Table("certificates").
		Select("certificates.*, forges.name AS forge_name, forges.image_url AS forge_image_url, swordsmiths.name AS swordsmith_name, swordsmiths.image_url AS swordsmith_image_url").
		Joins("LEFT JOIN forges ON forges.id = certificates.forge_id").
		Joins("LEFT JOIN swordsmiths ON swordsmiths.id = certificates.swordsmith_id")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(
			"certificates.number LIKE ? OR certificates.`no` LIKE ? OR forges.name LIKE ? OR swordsmiths.name LIKE ?",
			keyword, keyword, keyword, keyword,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询总数失败")
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("certificates.created_at DESC").Find(&certificates).Error; err != nil {
		return nil, errors.New("查询证书列表失败")
	}

	totalPage := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &models.CertificateListResponse{
		Certificates: certificates,
		Page:         params.Page,
		PageSize:     params.PageSize,
		Total:        total,
		TotalPage:    totalPage,
	}, nil
}

// GetCertificateByID 根据ID获取证书
func (s *CertificateService) GetCertificateByID(id uint) (*models.Certificate, error) {
	var certificate models.Certificate
	if err := s.db.Table("certificates").
		Select("certificates.*, forges.name AS forge_name, forges.image_url AS forge_image_url, swordsmiths.name AS swordsmith_name, swordsmiths.image_url AS swordsmith_image_url").
		Joins("LEFT JOIN forges ON forges.id = certificates.forge_id").
		Joins("LEFT JOIN swordsmiths ON swordsmiths.id = certificates.swordsmith_id").
		Where("certificates.id = ?", id).
		First(&certificate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("证书不存在")
		}
		return nil, errors.New("数据库查询错误")
	}
	return &certificate, nil
}

// GetCertificateByNo 根据证书序号获取证书
func (s *CertificateService) GetCertificateByNo(no string) (*models.Certificate, error) {
	trimmed := strings.TrimSpace(no)
	if trimmed == "" {
		return nil, errors.New("证书序号不能为空")
	}

	var certificate models.Certificate
	if err := s.db.Table("certificates").
		Select("certificates.*, forges.name AS forge_name, forges.image_url AS forge_image_url, swordsmiths.name AS swordsmith_name, swordsmiths.image_url AS swordsmith_image_url").
		Joins("LEFT JOIN forges ON forges.id = certificates.forge_id").
		Joins("LEFT JOIN swordsmiths ON swordsmiths.id = certificates.swordsmith_id").
		Where("certificates.`no` = ?", trimmed).
		First(&certificate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("证书不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	return &certificate, nil
}

// UpdateCertificate 更新证书
func (s *CertificateService) UpdateCertificate(id uint, req *models.UpdateCertificateRequest) (*models.Certificate, error) {
	certificate, err := s.GetCertificateByID(id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Number != "" {
		updates["number"] = req.Number
	}
	if req.OverallLength != "" {
		updates["overall_length"] = req.OverallLength
	}
	if req.BladeLength != "" {
		updates["blade_length"] = req.BladeLength
	}
	if req.HandleLength != "" {
		updates["handle_length"] = req.HandleLength
	}
	if req.BladeMaterial != "" {
		updates["blade_material"] = req.BladeMaterial
	}
	if req.BladeThickness != "" {
		updates["blade_thickness"] = req.BladeThickness
	}
	if req.KissanLength != "" {
		updates["kissan_length"] = req.KissanLength
	}
	if req.Hamon != "" {
		updates["hamon"] = req.Hamon
	}
	if req.Mekugi != nil {
		updates["mekugi"] = *req.Mekugi
	}
	if req.Sakihaba != "" {
		updates["sakihaba"] = req.Sakihaba
	}
	if req.Motohaba != "" {
		updates["motohaba"] = req.Motohaba
	}
	if req.SayaMaterial != "" {
		updates["saya_material"] = req.SayaMaterial
	}
	if req.TsubaMaterial != "" {
		updates["tsuba_material"] = req.TsubaMaterial
	}
	if req.HabakiMaterial != "" {
		updates["habaki_material"] = req.HabakiMaterial
	}
	if req.ItoSageoMaterial != "" {
		updates["ito_sageo_material"] = req.ItoSageoMaterial
	}
	if req.ForgeID != nil && *req.ForgeID > 0 {
		updates["forge_id"] = *req.ForgeID
	}
	if req.SwordsmithID != nil && *req.SwordsmithID > 0 {
		updates["swordsmith_id"] = *req.SwordsmithID
	}
	if req.DateCompleted != "" {
		updates["date_completed"] = req.DateCompleted
	}
	if req.No != "" {
		updates["no"] = req.No
	}

	if len(updates) == 0 {
		return certificate, nil
	}

	if err := s.db.Model(certificate).Updates(updates).Error; err != nil {
		return nil, errors.New("证书更新失败")
	}

	return s.GetCertificateByID(id)
}

// DeleteCertificate 删除证书
func (s *CertificateService) DeleteCertificate(id uint) error {
	certificate, err := s.GetCertificateByID(id)
	if err != nil {
		return err
	}

	if err := s.db.Delete(certificate).Error; err != nil {
		return errors.New("删除证书失败")
	}

	return nil
}
