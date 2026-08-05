package repo

 import (
        "context"
        "errors"

        "github.com/zyj/my-blog/internal/model"
  )

func CreateSession(ctx context.Context, session *model.Session) error {
	if db == nil {
		return errors.New("database is not initialized")
	}
	return db.WithContext(ctx).Create(session).Error
}

func RevokeSession(ctx context.Context, sessionID string) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).
	Where("session_id = ?", sessionID). 
	Delete(&model.Session{}). 
	Error
}

func IsSessionValidForUser(ctx context.Context, sessionID string, userID uint) (bool, error){
	 if db == nil {
                return false, errors.New("database is not initialized")
        }
		var count int64
		err := db.WithContext(ctx). 
		       Model(&model.Session{}). 
			   Where(
						"session_id = ? AND user_id = ? AND deleted_at IS NULL",
                        sessionID,
                        userID,
			   ). 
			   Count(&count). 
			   Error
		return count > 0, err
}

 func GetLatestSessionByUserID(
        ctx context.Context,
        userID uint,
  ) (model.Session, error) {
        if db == nil {
                return model.Session{}, errors.New("database is not initialized")
        }

        var session model.Session
        err := db.WithContext(ctx).
                Where("user_id = ?", userID).
                Order("created_at DESC").
                First(&session).
                Error

        return session, err
  }

  func RevokeSessionsByUserID(ctx context.Context, userID uint) error {
	if db == nil {
		return errors.New("database is not initialized")
	}
	return db.WithContext(ctx).
			Where("user_id = ?", userID).
			Delete(&model.Session{}).
			Error
  }
