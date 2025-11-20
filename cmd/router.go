package cmd

import (
	"net/http"

	"test-question/internal/infra"
	"test-question/internal/pkg/rpc/rpc_auth"
	"test-question/internal/pkg/timer"

	rpcQCreate "test-question/internal/rpc/question/create_question"
	rpcQDelete "test-question/internal/rpc/question/delete_question"
	rpcQGet "test-question/internal/rpc/question/get"
	rpcQList "test-question/internal/rpc/question/list"

	rpcACreate "test-question/internal/rpc/answer/create"
	rpcADelete "test-question/internal/rpc/answer/delete"
	rpcAGet "test-question/internal/rpc/answer/get"

	"test-question/internal/repository/answer"
	"test-question/internal/repository/question"
	"test-question/internal/repository/user"

	ucAuth "test-question/internal/usecase/auth"
	ucQCreate "test-question/internal/usecase/question/create"
	ucQDelete "test-question/internal/usecase/question/delete"
	ucQGet "test-question/internal/usecase/question/get_with_answers"
	ucQGetAll "test-question/internal/usecase/question/list"

	ucACreate "test-question/internal/usecase/answer/create"
	ucADelete "test-question/internal/usecase/answer/delete"
	ucAGet "test-question/internal/usecase/answer/get_by_id"

	"test-question/internal/pkg/uow"
)

func SetupRouter(resources *infra.Resources) http.Handler {
	// ==========================
	// Repositories
	// ==========================
	userRepo := user.NewRepository(resources.DB)
	questionRepo := question.NewRepository(resources.DB)
	answerRepo := answer.NewRepository(resources.DB)
	uowManager := uow.NewGormUoW(resources.DB)

	// ==========================
	// UseCases
	// ==========================
	authUseCase := ucAuth.NewUseCase(userRepo, resources.Logger)
	tm := timer.NewTimer()

	ucCreateQuestion := ucQCreate.NewUseCase(questionRepo, tm, resources.Logger)
	ucListQuestions := ucQGetAll.NewUseCase(questionRepo, resources.Logger)
	ucGetQuestion := ucQGet.NewUseCase(questionRepo, answerRepo, resources.Logger)
	ucDeleteQuestion := ucQDelete.NewUseCase(questionRepo, answerRepo, uowManager, resources.Logger)

	ucCreateAnswer := ucACreate.NewUseCase(answerRepo, questionRepo, tm, resources.Logger)
	ucDeleteAnswer := ucADelete.NewUseCase(answerRepo, resources.Logger)
	ucGetAnswer := ucAGet.NewUseCase(answerRepo, resources.Logger)

	// ==========================
	// HTTP Router (stdlib)
	// ==========================
	mux := http.NewServeMux()

	// --- Question handlers ---
	mux.Handle("POST /questions", rpcQCreate.NewHandler(ucCreateQuestion))
	mux.Handle("GET /questions", rpcQList.NewHandler(ucListQuestions))
	mux.Handle("GET /questions/{id}", rpcQGet.NewHandler(ucGetQuestion))
	mux.Handle("DELETE /questions/{id}", rpcQDelete.NewHandler(ucDeleteQuestion))

	// --- Answer handlers ---
	mux.Handle("POST /questions/{id}/answers", rpcACreate.NewHandler(ucCreateAnswer))
	mux.Handle("GET /answers/{id}", rpcAGet.NewHandler(ucGetAnswer))
	mux.Handle("DELETE /answers/{id}", rpcADelete.NewHandler(ucDeleteAnswer))

	// ==========================
	// Wrap with middleware
	// ==========================
	handler := rpc_auth.BasicAuthMiddleware(authUseCase)(mux)

	return handler
}
