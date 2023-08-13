# log

A simple package that wraps [zap logger](https://github.com/uber-go/zap) with [lumberjack](https://github.com/natefinch/lumberjack). This package provides a simpler way to init and use zap logger. 

### Installation

`go get -u github.com/olee12/log`



### Example

````go
type User struct {
	FirstName string
	LastName  string
	Age       int
	Address   Address
	Phone     string
}

type Address struct {
	Street   string
	PostCode string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomUser() User {
	return User{
		FirstName: RandStringRunes(5),
		LastName:  RandStringRunes(5),
		Age:       rand.Intn(30) + 10,
		Address: Address{
			Street:   RandStringRunes(20),
			PostCode: RandStringRunes(6),
		},
		Phone: RandStringRunes(10),
	}
}

func main() {
	packageLevelLoggerExample()
	childLoggerExample(100000)
}

func packageLevelLoggerExample() {
	packageDebugLogger()
	packageErrorLogger()
	packageInfoLogger()
	packageWarnLogger()
}

func childLoggerExample(N int) {
	log.Configure(*log.DefaultRotateLoggerConfig)
	logLevels := []log.Level{log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.PanicLevel, log.FatalLevel, log.DPanicLevel}
	ticker := time.NewTicker(2 * time.Second)

	// change the log level every concurrently
	go func() {
		for range ticker.C {
			log.SetLevel(logLevels[rand.Intn(len(logLevels))])
		}
	}()

	for i := 0; i < N; i++ {
		logger := log.WithField("request_id", RandStringRunes(15))
		joiner := make(chan struct{}, 4)
		go func() {
			tryDebugLogger(logger)
			joiner <- struct{}{}
		}()

		go func() {
			tryInfoLogger(logger)
			joiner <- struct{}{}
		}()

		go func() {
			tryWarnLogger(logger)
			joiner <- struct{}{}
		}()

		go func() {
			tryErrorLogger(logger)
			joiner <- struct{}{}
		}()

		// join
		for a := 0; a < 4; a++ {
			<-joiner
		}

	}
	ticker.Stop()
}

func tryInfoLogger(logger *log.LogEntry) {
	logger.InfoWith("[with] message with some user info", log.Fields{"user1": RandomUser(), "user2": RandomUser()})
	logger.InfoWith("[with] only message", log.Fields{})
	logger.Info("[] simple log")
	logger.Infov("[v] with user", zap.Any("user", RandomUser()))
	logger.Infow("[w] with user", "user", RandomUser())
	logger.Infof("[f]user: %v", RandomUser())
}

func tryDebugLogger(logger *log.LogEntry) {
	logger.DebugWith("[with]message with some user info", log.Fields{"user1": RandomUser(), "user2": RandomUser()})
	logger.DebugWith("[with]only message", log.Fields{})
	logger.Debug("[] simple log")
	logger.Debugv("[v] with user", zap.Any("user", RandomUser()))
	logger.Debugw("[w] with user", "user", RandomUser())
	logger.Debugf("[f] user: %v", RandomUser())
}

func tryErrorLogger(logger *log.LogEntry) {
	logger.ErrorWith("[with] message with some user info", log.Fields{"user1": RandomUser(), "user2": RandomUser()})
	logger.ErrorWith("[with] only message", log.Fields{})
	logger.Error("[] simple log")
	logger.Errorv("[v] with user", zap.Any("user", RandomUser()))
	logger.Errorw("[w] with user", "user", RandomUser())
	logger.Errorf("[f] user: %v", RandomUser())
}

func tryWarnLogger(logger *log.LogEntry) {
	logger.WarnWith("[with] message with some user info", log.Fields{"user1": RandomUser(), "user2": RandomUser()})
	logger.WarnWith("[with] only message", log.Fields{})
	logger.Warn("[] simple log")
	logger.Warnv("[v] with user", zap.Any("user", RandomUser()))
	logger.Warnw("[w] with user", "user", RandomUser())
	logger.Warnf("[f] user: %v", RandomUser())
}

func packageInfoLogger() {
	log.InfoWith("[with] message with some user info", log.Fields{"user1": RandomUser(), "user2": RandomUser()})
	log.InfoWith("[with] only message", log.Fields{})
	log.Info("[] simple log")
	log.Infov("[v] with user", zap.Any("user", RandomUser()))
	log.Infow("[w] with user", "user", RandomUser())
	log.Infof("[f]user: %v", RandomUser())
}

func packageDebugLogger() {
	log.DebugWith("[with]message with some user info", log.Fields{"user1": RandomUser(), "user2": RandomUser()})
	log.DebugWith("[with]only message", log.Fields{})
	log.Debug("[] simple log")
	log.Debugv("[v] with user", zap.Any("user", RandomUser()))
	log.Debugw("[w] with user", "user", RandomUser())
	log.Debugf("[f] user: %v", RandomUser())
}

func packageErrorLogger() {
	log.ErrorWith("[with] message with some user info", log.Fields{"user1": RandomUser(), "user2": RandomUser()})
	log.ErrorWith("[with] only message", log.Fields{})
	log.Error("[] simple log")
	log.Errorv("[v] with user", zap.Any("user", RandomUser()))
	log.Errorw("[w] with user", "user", RandomUser())
	log.Errorf("[f] user: %v", RandomUser())
}

func packageWarnLogger() {
	log.WarnWith("[with] message with some user info", log.Fields{"user1": RandomUser(), "user2": RandomUser()})
	log.WarnWith("[with] only message", log.Fields{})
	log.Warn("[] simple log")
	log.Warnv("[v] with user", zap.Any("user", RandomUser()))
	log.Warnw("[w] with user", "user", RandomUser())
	log.Warnf("[f] user: %v", RandomUser())
}
````