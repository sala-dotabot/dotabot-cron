package repository

type MapRespository struct {
	holder map[TelegramMatchSubscription]uint64
}

func CreateMapRepository() SubscriptionRepository {
	return &MapRespository{holder: make(map[TelegramMatchSubscription]uint64)}
}

func (this MapRespository) FindAll() ([]TelegramMatchSubscription, error) {
	var keys []TelegramMatchSubscription
	for k := range this.holder {
		keys = append(keys, k)
	}
	return keys, nil
}

func (this MapRespository) FindByChatId(chatId int64) (subscriptions []TelegramMatchSubscription, err error) {
	for k := range this.holder {
		if k.ChatId == chatId {
			subscriptions = append(subscriptions, k)
		}
	}
	return
}

func (this MapRespository) GetLastKnownMatchId(subscription TelegramMatchSubscription) (int64, error) {
	matchId, ok := this.holder[subscription]
	if ok {
		return int64(matchId), nil
	} else {
		return -1, nil
	}
}

func (this MapRespository) SaveLastKnownMatchId(subscription TelegramMatchSubscription, matchId uint64) error {
	this.holder[subscription] = matchId
	return nil
}

func (this MapRespository) RemoveLastKnownMatchId(subscription TelegramMatchSubscription) error {
	delete(this.holder, subscription)
	return nil
}
