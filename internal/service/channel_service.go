package service

import (
	"errors"
	"netradio/internal/controller/requests"
	"netradio/internal/controller/responses"
	"netradio/internal/model"
	"netradio/internal/repository"
	webrtchelper "netradio/pkg/webrtc"
	"strconv"
)

type ChannelService interface {
	GetChannels() (responses.GetChannelsResponse, error)
	GetChannel(r requests.GetChannelRequest) (responses.GetChannelResponse, error)
	CreateChannel(r requests.CreateChannelRequest) error
	DeleteChannel(r requests.DeleteChannelRequest) error
	UpdateChannel(r requests.UpdateChannelRequest) (responses.UpdateChannelResponse, error)
	StartChannel(r requests.StartChannelRequest) error
	StopChannel(r requests.StopChannelRequest) error
	UploadLogo(r requests.UploadLogoRequest) error
	AddTrack(r requests.AddTrackRequest) error
	GetCurrentTrack(r requests.GetCurrentTrackRequest) (responses.GetCurrentTrackResponse, error)
}

func NewChannelService() ChannelService {
	return &ChannelServiceImpl{
		db: repository.NewChannelDB(),
	}
}

type ChannelServiceImpl struct {
	db repository.ChannelDB
}

func (s *ChannelServiceImpl) GetChannels() (responses.GetChannelsResponse, error) {
	var res responses.GetChannelsResponse
	channels, err := s.db.GetChannels()
	if err != nil {
		return res, err
	}

	res.Channels = channels

	return res, nil
}

func (s *ChannelServiceImpl) GetChannel(r requests.GetChannelRequest) (responses.GetChannelResponse, error) {
	var res responses.GetChannelResponse
	channel, err := s.db.GetChannelById(r.ID)
	if err != nil {
		return res, err
	}
	if channel == nil {
		res.Found = false
		return res, err
	}

	res.Found = true
	res.ID = strconv.Itoa(channel.ID)
	res.Title = channel.Title
	res.Description = channel.Description
	res.Status = int(channel.Status)
	res.Logo = channel.Logo

	return res, nil
}

func (s *ChannelServiceImpl) CreateChannel(r requests.CreateChannelRequest) error {
	var channel model.ChannelInfo
	channel.Title = r.Title
	channel.Description = r.Description
	channel.Status = model.StoppedChannel
	return s.db.CreateChannel(channel)
}

func (s *ChannelServiceImpl) DeleteChannel(r requests.DeleteChannelRequest) error {
	return s.db.DeleteChannel(r.ID)
}

func (s *ChannelServiceImpl) UpdateChannel(r requests.UpdateChannelRequest) (responses.UpdateChannelResponse, error) {
	var res responses.UpdateChannelResponse
	channel, err := s.db.GetChannelById(r.ID)
	if err != nil {
		return res, err
	}
	if channel == nil {
		res.Found = false
		return res, nil
	}

	res.Found = true
	if r.Title != nil {
		channel.Title = *r.Title
	}
	if r.Description != nil {
		channel.Description = *r.Description
	}

	err = s.db.UpdateChannel(*channel)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *ChannelServiceImpl) StartChannel(r requests.StartChannelRequest) error {
	return s.db.ChangeChannelStatus(r.ID, model.ActiveChannel)
}

func (s *ChannelServiceImpl) StopChannel(r requests.StopChannelRequest) error {
	return s.db.ChangeChannelStatus(r.ID, model.StoppedChannel)
}

func (s *ChannelServiceImpl) UploadLogo(r requests.UploadLogoRequest) error {
	return s.db.ChangeLogo(r.ID, r.Logo)
}

func (s *ChannelServiceImpl) AddTrack(r requests.AddTrackRequest) error {
	track, err := repository.NewTrackDB().GetTrackById(r.TrackID)
	if err != nil {
		return err
	}

	return s.db.AddTrackToSchedule(r.ID, r.TrackID, r.StartDate, r.StartDate.Add(track.Duration))
}

func (s *ChannelServiceImpl) GetCurrentTrack(r requests.GetCurrentTrackRequest) (responses.GetCurrentTrackResponse, error) {
	var res responses.GetCurrentTrackResponse
	track, err := s.db.GetCurrentTrack(r.ID)
	if err != nil {
		return res, err
	}
	if track == nil {
		return res, errors.New("Track not found")
	}

	res.ID = strconv.Itoa(track.ID)
	res.Title = track.Title
	res.Perfomancer = track.Perfomancer
	res.Year = track.Year
	res.Duration = track.Duration
	res.CurrentTime = webrtchelper.GetCurrentTrackTime(r.ID)
	likeCount, err := repository.NewTrackDB().LikeCount(res.ID)
	if err != nil {
		return res, err
	}
	res.LikeCount = likeCount
	if r.UserID != nil {
		liked, err := repository.NewTrackDB().IsTrackLiked(r.ID, *r.UserID)
		if err != nil {
			return res, err
		}
		res.Liked = liked
	} else {
		res.Liked = false
	}

	return res, nil
}
