package services

import (
	"errors"
	"steam-backend/config"
	"steam-backend/logger"
	"steam-backend/models"
	"strconv"

	"gorm.io/gorm"
)

// RecordService 兑换记录服务结构体
type RecordService struct {
	db                    *gorm.DB
	activationCodeService *ActivationCodeService
}

// NewRecordService 创建新的兑换记录服务实例
func NewRecordService() *RecordService {
	return &RecordService{
		db:                    config.GetDB(),
		activationCodeService: NewActivationCodeService(),
	}
}

// CreateExchangeRecord 创建兑换记录
func (s *RecordService) CreateExchangeRecord(req *models.CreateExchangeRecordRequest) (*models.ExchangeRecord, error) {
	// 校验激活码格式（快速失败，避免无效请求进入事务）
	ok := s.activationCodeService.codeGenerator.ValidateActivationCodeFormat(req.Code)
	if !ok {
		return nil, errors.New("激活码格式不正确")
	}

	// 用户两次好友码验证
	steamId, err := s.convertFriendCodeToSteamID(req.FriendCode, req.RepeatFriendCode)
	if err != nil {
		return nil, err
	}

	// 使用数据库事务确保数据一致性
	var record models.ExchangeRecord
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 使用排他锁查询激活码，防止并发问题
		var activationCode models.ActivationCode
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("code = ?", req.Code).
			First(&activationCode).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("未查询到激活码")
			}
			return errors.New("查询激活码失败")
		}

		// 检查激活码状态
		if activationCode.Status != models.StatusNormalOfActivationCode {
			return errors.New("激活码已被使用")
		}

		// 更新激活码状态为已使用
		if err := tx.Model(&activationCode).Update("status", models.StatusUsedOfActivationCode).Error; err != nil {
			return errors.New("更新激活码状态失败")
		}

		// 创建兑换记录
		record = models.ExchangeRecord{
			Code:    req.Code,
			SteamID: steamId,
			Points:  activationCode.Points,                // 从激活码获取点数
			Status:  models.StatusPendingOfExchangeRecord, // 默认状态为待兑换
		}

		// 保存兑换记录
		if err := tx.Create(&record).Error; err != nil {
			return errors.New("创建兑换记录失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Info("成功创建兑换记录: ID=%d, SteamID=%d, 激活码=%s, 点数=%d",
		record.ID, record.SteamID, record.Code, record.Points)

	return &record, nil
}

// GetExchangeRecords 获取兑换记录列表（带分页和搜索）
func (s *RecordService) GetExchangeRecords(params *models.ExchangeRecordListRequest) (*models.ExchangeRecordListResponse, error) {
	var records []models.ExchangeRecord
	var total int64

	// 构建查询
	query := s.db.Model(&models.ExchangeRecord{})

	// 添加激活码筛选
	if params.Code != "" {
		query = query.Where("code = ?", params.Code)
	}

	// 添加SteamID筛选
	if params.SteamID != 0 {
		query = query.Where("steam_id = ?", params.SteamID)
	}

	// 添加状态筛选
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询总数失败")
	}

	// 设置默认分页参数
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// 默认预加载礼物记录
	query = query.Preload("GiftRecords")

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, errors.New("查询兑换记录列表失败")
	}

	// 计算总页数
	totalPage := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &models.ExchangeRecordListResponse{
		Records:   records,
		Page:      params.Page,
		PageSize:  params.PageSize,
		Total:     total,
		TotalPage: totalPage,
	}, nil
}

// GetExchangeRecordByID 根据ID获取兑换记录信息
func (s *RecordService) GetExchangeRecordByID(recordID uint) (*models.ExchangeRecord, error) {
	return s.GetExchangeRecordByIDWithOptions(recordID, true)
}

// GetExchangeRecordByIDWithOptions 根据ID获取兑换记录信息（支持选项）
func (s *RecordService) GetExchangeRecordByIDWithOptions(recordID uint, includeGifts bool) (*models.ExchangeRecord, error) {
	var record models.ExchangeRecord

	query := s.db
	if includeGifts {
		query = query.Preload("GiftRecords")
	}

	if err := query.First(&record, recordID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("兑换记录不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	return &record, nil
}

// getSteamIDFromLink 根据提供的好友链接获取 STEAMID
func (s *RecordService) convertFriendCodeToSteamID(friendCode, repeatFriendCode string) (uint64, error) {
	if friendCode != repeatFriendCode {
		return 0, errors.New("两次输入的好友码不一致")
	}

	// 将字符串转换为uint64
	codeValue, err := strconv.ParseUint(friendCode, 10, 64)
	if err != nil {
		return 0, errors.New("好友码格式有误")
	}

	// SteamID64的基础值
	const steamId64Base = 76561197960265728

	// 计算最终的SteamID64
	steamID64 := steamId64Base + codeValue

	return steamID64, nil
}
