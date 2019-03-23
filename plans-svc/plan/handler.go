package plan

import (
	"net/http"
	"time"

	goKitLog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	lib "github.com/syedomair/plan-api/lib"
	"github.com/syedomair/plan-api/models"
)

type PlanEnv struct {
	Logger   goKitLog.Logger
	PlanRepo PlanRepositoryInterface
	Common   lib.CommonService
}

func (env *PlanEnv) CreatePlan(w http.ResponseWriter, r *http.Request) {

	env.Logger.Log("METHOD", "CreatePlan", "SPOT", "method start")
	start := time.Now()
	_, err := env.Common.GetUserClientFromToken(r)

	if err != nil {
		env.Common.ErrorResponseHelper(w, "5011", err.Error(), http.StatusBadRequest)
		return
	}

	var pathParamConf map[string]string
	pathParamConf = make(map[string]string)

	var paramConf map[string]models.ParamConf
	paramConf = make(map[string]models.ParamConf)
	paramConf["title"] = models.ParamConf{Required: true, Type: lib.STRING_SMALL, EmptyAllowed: false}
	paramConf["status"] = models.ParamConf{Required: false, Type: lib.INT, EmptyAllowed: false}
	paramConf["cost"] = models.ParamConf{Required: true, Type: lib.INT, EmptyAllowed: false}
	paramConf["validity"] = models.ParamConf{Required: true, Type: lib.INT, EmptyAllowed: false}

	_, paramMap, errCode, err := env.Common.ValidateInputParameters(r, paramConf, pathParamConf)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5"+errCode, err.Error(), http.StatusBadRequest)
		return
	}

	planId, err := env.PlanRepo.Create(paramMap)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5019", lib.ERROR_UNEXPECTED, http.StatusInternalServerError)
	}

	responsePlanId := map[string]string{"plan_id": planId}

	env.Logger.Log("METHOD", "CreatePlan", "SPOT", "METHOD END", "time_spent", time.Since(start))
	env.Common.SuccessResponseHelper(w, responsePlanId, http.StatusCreated)
}

func (env *PlanEnv) GetAll(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	env.Logger.Log("METHOD", "GetAll", "SPOT", "method start", "time_start", start)

	_, err := env.Common.GetUserClientFromToken(r)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5020", err.Error(), http.StatusBadRequest)
		return
	}

	limit, offset, orderby, sort, err := env.Common.ValidateQueryString(r, "100", "0", "created_at", "asc")
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5021", err.Error(), http.StatusBadRequest)
		return
	}
	plans, count, err := env.PlanRepo.GetAll(limit, offset, orderby, sort)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5022", lib.ERROR_UNEXPECTED, http.StatusInternalServerError)
	}
	env.Logger.Log("METHOD", "GetAll", "SPOT", "method end", "time_spent", time.Since(start))
	env.Common.SuccessResponseList(w, plans, offset, limit, count)
}
func (env *PlanEnv) Get(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	env.Logger.Log("METHOD", "Get", "SPOT", "method start", "time_start", start)
	_, err := env.Common.GetUserClientFromToken(r)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5023", err.Error(), http.StatusBadRequest)
		return
	}
	pathParams := mux.Vars(r)
	planId := pathParams["plan_id"]

	if err := env.Common.ValidateId(planId, "plan_id"); err != nil {
		env.Common.ErrorResponseHelper(w, "5024", err.Error(), http.StatusBadRequest)
		return
	}
	plan, err := env.PlanRepo.Get(planId)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5025", lib.ERROR_UNEXPECTED, http.StatusInternalServerError)
		return
	}
	env.Logger.Log("METHOD", "Get", "SPOT", "method end", "time_spent", time.Since(start))
	env.Common.SuccessResponseHelper(w, plan, http.StatusOK)
}

func (env *PlanEnv) Update(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	env.Logger.Log("METHOD", "Update", "SPOT", "method start", "time_start", start)
	_, err := env.Common.GetUserClientFromToken(r)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5026", err.Error(), http.StatusBadRequest)
		return
	}

	var pathParamConf map[string]string
	pathParamConf = make(map[string]string)
	pathParamConf["plan_id"] = ""

	var paramConf map[string]models.ParamConf
	paramConf = make(map[string]models.ParamConf)
	paramConf["title"] = models.ParamConf{Required: false, Type: lib.STRING_SMALL, EmptyAllowed: false}
	paramConf["status"] = models.ParamConf{Required: false, Type: lib.INT, EmptyAllowed: false}
	paramConf["cost"] = models.ParamConf{Required: false, Type: lib.INT, EmptyAllowed: false}
	paramConf["validity"] = models.ParamConf{Required: false, Type: lib.INT, EmptyAllowed: false}

	pathParamValue, paramMap, errCode, err := env.Common.ValidateInputParameters(r, paramConf, pathParamConf)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5"+errCode, err.Error(), http.StatusBadRequest)
		return
	}

	planId := pathParamValue["plan_id"]

	err = env.PlanRepo.Update(paramMap, planId)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5034", lib.ERROR_UNEXPECTED, http.StatusInternalServerError)
		return
	}
	responsePlanId := map[string]string{"plan_id": planId}
	env.Logger.Log("METHOD", "Update", "SPOT", "method end", "time_spent", time.Since(start))
	env.Common.SuccessResponseHelper(w, responsePlanId, http.StatusOK)
}

func (env *PlanEnv) Delete(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	env.Logger.Log("METHOD", "Delete", "SPOT", "method start", "time_start", start)

	_, err := env.Common.GetUserClientFromToken(r)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5035", err.Error(), http.StatusBadRequest)
		return
	}
	pathParams := mux.Vars(r)
	planId := pathParams["plan_id"]
	if err := env.Common.ValidateId(planId, "plan_id"); err != nil {
		env.Common.ErrorResponseHelper(w, "5036", err.Error(), http.StatusBadRequest)
		return
	}

	plan := models.Plan{}
	plan.Id = planId

	err = env.PlanRepo.Delete(plan)
	if err != nil {
		env.Common.ErrorResponseHelper(w, "5037", lib.ERROR_UNEXPECTED, http.StatusInternalServerError)
		return
	}

	responsePlanId := map[string]string{"plan_id": planId}
	env.Logger.Log("METHOD", "Delete", "SPOT", "method end", "time_spent", time.Since(start))
	env.Common.SuccessResponseHelper(w, responsePlanId, http.StatusOK)
}
