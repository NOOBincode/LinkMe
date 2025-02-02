// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"github.com/GoSimplicity/LinkMe/internal/api"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/check"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/email"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/es"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/post"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/publish"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/sms"
	"github.com/GoSimplicity/LinkMe/internal/job"
	"github.com/GoSimplicity/LinkMe/internal/mock"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/GoSimplicity/LinkMe/utils/jwt"
)

import (
	_ "github.com/google/wire"
)

// Injectors from wire.go:

func InitWebServer() *Cmd {
	db := InitDB()
	node := InitializeSnowflakeNode()
	logger := InitLogger()
	enforcer := InitCasbin(db)
	userDAO := dao.NewUserDAO(db, node, logger, enforcer)
	cmdable := InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache, logger)
	typedClient := InitES()
	searchDAO := dao.NewSearchDAO(db, typedClient, logger)
	searchRepository := repository.NewSearchRepository(searchDAO)
	userService := service.NewUserService(userRepository, logger, searchRepository)
	handler := jwt.NewJWTHandler(cmdable)
	client := InitSaramaClient()
	syncProducer := InitSyncProducer(client)
	producer := sms.NewSaramaSyncProducer(syncProducer, logger)
	emailProducer := email.NewSaramaSyncProducer(syncProducer, logger)
	userHandler := api.NewUserHandler(userService, handler, producer, emailProducer, enforcer)
	postDAO := dao.NewPostDAO(db, logger)
	postCache := cache.NewPostCache(cmdable)
	asynqClient := InitAsynqClient()
	postRepository := repository.NewPostRepository(postDAO, logger, postCache, asynqClient)
	postProducer := post.NewSaramaSyncProducer(syncProducer)
	checkProducer := check.NewSaramaCheckProducer(syncProducer)
	interactiveDAO := dao.NewInteractiveDAO(db, logger)
	interactiveRepository := repository.NewInteractiveRepository(interactiveDAO, logger)
	postService := service.NewPostService(postRepository, logger, postProducer, checkProducer, interactiveRepository)
	interactiveService := service.NewInteractiveService(interactiveRepository, logger)
	postHandler := api.NewPostHandler(postService, interactiveService)
	historyCache := cache.NewHistoryCache(logger, cmdable)
	historyRepository := repository.NewHistoryRepository(logger, historyCache)
	historyService := service.NewHistoryService(historyRepository, logger)
	historyHandler := api.NewHistoryHandler(historyService)
	checkDAO := dao.NewCheckDAO(db, logger)
	checkRepository := repository.NewCheckRepository(checkDAO, logger)
	activityDAO := dao.NewActivityDAO(db, logger)
	activityRepository := repository.NewActivityRepository(activityDAO)
	publishProducer := publish.NewSaramaSyncProducer(syncProducer, logger)
	checkService := service.NewCheckService(checkRepository, searchRepository, logger, activityRepository, publishProducer)
	checkHandler := api.NewCheckHandler(checkService)
	v := InitMiddlewares(handler, logger)
	apiDAO := dao.NewApiDAO(db, logger)
	permissionDAO := dao.NewPermissionDAO(db, logger, enforcer, apiDAO)
	permissionRepository := repository.NewPermissionRepository(logger, permissionDAO)
	permissionService := service.NewPermissionService(logger, permissionRepository)
	permissionHandler := api.NewPermissionHandler(permissionService, logger)
	rankingRedisCache := cache.NewRankingRedisCache(cmdable, logger)
	rankingLocalCache := cache.NewRankingLocalCache(logger)
	rankingRepository := repository.NewRankingCache(rankingRedisCache, rankingLocalCache, logger)
	rankingService := service.NewRankingService(interactiveRepository, postRepository, rankingRepository, logger)
	rankingHandler := api.NewRakingHandler(rankingService)
	plateDAO := dao.NewPlateDAO(logger, db)
	plateRepository := repository.NewPlateRepository(logger, plateDAO)
	plateService := service.NewPlateService(logger, plateRepository)
	plateHandler := api.NewPlateHandler(plateService, enforcer)
	activityService := service.NewActivityService(activityRepository)
	activityHandler := api.NewActivityHandler(activityService, enforcer)
	commentDAO := dao.NewCommentDAO(db, logger)
	commentCache := cache.NewCommentCache(cmdable)
	commentRepository := repository.NewCommentRepository(commentDAO,commentCache)
	commentService := service.NewCommentService(commentRepository,checkProducer)
	commentHandler := api.NewCommentHandler(commentService)
	searchService := service.NewSearchService(searchRepository)
	searchHandler := api.NewSearchHandler(searchService)
	relationDAO := dao.NewRelationDAO(db, logger)
	relationCache := cache.NewRelationCache(cmdable)
	relationRepository := repository.NewRelationRepository(relationDAO, relationCache, logger)
	relationService := service.NewRelationService(relationRepository)
	relationHandler := api.NewRelationHandler(relationService)
	lotteryDrawDAO := dao.NewLotteryDrawDAO(db, logger)
	lotteryDrawRepository := repository.NewLotteryDrawRepository(lotteryDrawDAO, logger)
	lotteryDrawService := service.NewLotteryDrawService(lotteryDrawRepository, logger)
	lotteryDrawHandler := api.NewLotteryDrawHandler(lotteryDrawService)
	roleDAO := dao.NewRoleDAO(db, logger, enforcer, permissionDAO)
	roleRepository := repository.NewRoleRepository(logger, roleDAO)
	roleService := service.NewRoleService(roleRepository, permissionRepository, logger)
	menuDAO := dao.NewMenuDAO(db, logger)
	menuRepository := repository.NewMenuRepository(logger, menuDAO)
	menuService := service.NewMenuService(logger, menuRepository)
	apiRepository := repository.NewApiRepository(logger, apiDAO)
	apiService := service.NewApiService(logger, apiRepository)
	roleHandler := api.NewRoleHandler(roleService, menuService, apiService, permissionService, logger)
	menuHandler := api.NewMenuHandler(menuService, logger)
	apiHandler := api.NewApiHandler(apiService, logger)
	engine := InitWeb(userHandler, postHandler, historyHandler, checkHandler, v, permissionHandler, rankingHandler, plateHandler, activityHandler, commentHandler, searchHandler, relationHandler, lotteryDrawHandler, roleHandler, menuHandler, apiHandler)
	eventConsumer := post.NewEventConsumer(interactiveRepository, historyRepository, client, syncProducer, logger)
	smsDAO := dao.NewSmsDAO(db, logger)
	smsCache := cache.NewSMSCache(cmdable)
	tencentSms := InitSms()
	smsRepository := repository.NewSmsRepository(smsDAO, smsCache, logger, tencentSms)
	smsConsumer := sms.NewSMSConsumer(smsRepository, client, logger, smsCache)
	emailCache := cache.NewEmailCache(cmdable)
	emailRepository := repository.NewEmailRepository(emailCache, logger)
	emailConsumer := email.NewEmailConsumer(emailRepository, client, logger)
	publishPostEventConsumer := publish.NewPublishPostEventConsumer(postRepository, client, syncProducer, logger)
	esConsumer := es.NewEsConsumer(client, logger, searchRepository)
	checkEventConsumer := check.NewCheckEventConsumer(checkRepository, client, syncProducer, logger)
	postDeadLetterConsumer := post.NewPostDeadLetterConsumer(interactiveRepository, historyRepository, client, logger)
	publishDeadLetterConsumer := publish.NewPublishDeadLetterConsumer(postRepository, client, logger)
	checkDeadLetterConsumer := check.NewCheckDeadLetterConsumer(checkRepository, client, logger)
	v2 := InitConsumers(eventConsumer, smsConsumer, emailConsumer, publishPostEventConsumer, esConsumer, checkEventConsumer, postDeadLetterConsumer, publishDeadLetterConsumer, checkDeadLetterConsumer)
	mockUserRepository := mock.NewMockUserRepository(db, logger, enforcer)
	refreshCacheTask := job.NewRefreshCacheTask(postCache, logger)
	timedTask := job.NewTimedTask(logger)
	routes := job.NewRoutes(refreshCacheTask, timedTask)
	server := InitAsynqServer()
	scheduler := InitScheduler()
	timedScheduler := job.NewTimedScheduler(scheduler)
	cmd := &Cmd{
		Server:    engine,
		Consumer:  v2,
		Mock:      mockUserRepository,
		Routes:    routes,
		Asynq:     server,
		Scheduler: timedScheduler,
	}
	return cmd
}
