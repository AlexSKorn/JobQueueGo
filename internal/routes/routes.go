package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Job struct {
	ID     int    `json:"ID"`
	Type   string `json:"Type"`
	Status string `json:"Status"`
}

type JobQueue struct {
	Jobs     []*Job
	JobMutex sync.Mutex
}

const (
	QUEUED      string = "QUEUED"
	IN_PROGRESS string = "IN_PROGRESS"
	CONCLUDED   string = "CONCLUDED"
	DEQUEUED    string = "DEQUEUED"

	TIME_CRITICAL     string = "TIME_CRITICAL"
	NOT_TIME_CRITICAL string = "NOT_TIME_CRITICAL"
)

func NewRouter(jq *JobQueue) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/jobs/enqueue", jq.enqueueJob).Methods("POST")
	router.HandleFunc("/jobs/dequeue", jq.dequeueJob).Methods("GET")
	router.HandleFunc("/jobs/{job_id}", jq.getJob).Methods("GET")
	router.HandleFunc("/jobs/{job_id}/conclude", jq.concludeJob).Methods("PUT")
	router.Handle("/", router)

	return router
}

func (jq *JobQueue) enqueueJob(w http.ResponseWriter, r *http.Request) {
	var job Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//In the future I would check if it was not equal to the available options
	switch job.Type {
	case TIME_CRITICAL, NOT_TIME_CRITICAL:
		// valid job type, do nothing
	default:
		http.Error(w, "Invalid job type", http.StatusBadRequest)
		return
	}

	switch job.Status {
	case QUEUED, IN_PROGRESS, CONCLUDED:
		// valid job type, do nothing
	default:
		http.Error(w, "Invalid job status", http.StatusBadRequest)
		return
	}

	job.ID = len(jq.Jobs) + 1
	job.Status = QUEUED
	jq.JobMutex.Lock()
	jq.Jobs = append(jq.Jobs, &job)
	jq.JobMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job.ID)
	// fmt.Fprintln(w, "Enqueue Job")
}

// If we wanted to reference it later we could add it to a database in a prod case or in this local case
// another data structure.
func (jq *JobQueue) dequeueJob(w http.ResponseWriter, r *http.Request) {
	var job *Job
	jq.JobMutex.Lock()
	for _, j := range jq.Jobs {
		if j.Status == QUEUED || j.Status == IN_PROGRESS {
			job = j
			job.Status = DEQUEUED
			// Seems like we don't need to remove the job and we want to keep it for historical reasons?
			// Remove the job from the queue.
			//jq.Jobs = append(jq.Jobs[:i], jq.Jobs[i+1:]...)
			break
		}
	}
	jq.JobMutex.Unlock()
	if job == nil {
		//400 or 404 no sure
		http.Error(w, "No jobs available", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(job)
	// fmt.Fprintln(w, "Dequeue Job")
}

func (jq *JobQueue) getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := strconv.Atoi(vars["job_id"])

	if err != nil {
		http.Error(w, "Invalid Job ID", 400)
		return
	}

	var job *Job
	jq.JobMutex.Lock()
	for _, j := range jq.Jobs {
		if j.ID == jobID {
			job = j
			break
		}
	}
	jq.JobMutex.Unlock()

	if job == nil {
		http.Error(w, "Job not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(job)
	// fmt.Fprintln(w, "Get Job")
}

func (jq *JobQueue) concludeJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	jobID, err := strconv.Atoi(vars["job_id"])
	if err != nil {
		http.Error(w, "Invalid Job ID", 400)
		return
	}
	jq.JobMutex.Lock()
	defer jq.JobMutex.Unlock()

	var job *Job
	for _, j := range jq.Jobs {
		if j.ID == jobID {
			job = j
			if j.Status == CONCLUDED || j.Status == DEQUEUED {
				http.Error(w, "job already concluded/dequeued", 400)
				return
			}
			job.Status = CONCLUDED
		}
	}

	if job == nil {
		http.Error(w, "Job not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(job)
	// fmt.Fprintln(w, "Conclude")
}
