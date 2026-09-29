/* SPDX-License-Identifier: MPL-2.0 */
#include "workload.h"
#include "util.h"
#include "log.h"
#include <errno.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <time.h>
#include <unistd.h>

_Static_assert(ATOMIC_INT_LOCK_FREE == 2, "workloads require lock-free signal flags");
static _Atomic sig_atomic_t stopping;
static void stop_handler(int signal_number) { (void)signal_number; stopping=1; }

int workload_run(const runner_opts_t *opts,int plan_only) {
    wl_config_t config={0}; wl_catalog_t catalog={0}; wl_environment_t env={0};
    wl_metrics_t *metrics=NULL; ulab_error_t err={0},cleanup_err={0};
    char run_dir[ULAB_MAX_PATH]={0},run_id[ULAB_MAX_ID],reason[ULAB_MAX_ERR]={0};
    int rc=ULAB_OK,cleanup_rc=ULAB_OK,signals=0,reported=0;
    struct sigaction old_int,old_term,action;
    if(wl_config_load(opts->scenario_path,opts,&config,&err) || wl_catalog_load(&config,&catalog,&err)) { rc=ULAB_ESCENARIO;goto done; }
    if(plan_only) {
        printf("WORKLOAD %s\n",config.name);
        printf("world maxima: sites=%u live_ues_per_site=%u backend_sims=%u\n",
            config.environment->world.sites_per_network,config.environment->world.ues_per_site,config.environment->world.sims_per_network);
        printf("initial: sites=%zu provisioned=%zu attached=%zu\n",config.initial.sites,config.initial.provisioned,config.initial.attached);
        for(size_t i=0;i<config.phase_count;i++) {
            wl_phase_t *p=&config.phases[i];
            printf("%-24s sites=%zu provisioned=%zu attached=%zu active=%zu sessions=%zu rps=%g ramp=%gs reach=%gs hold=%gs\n",
                p->name,p->target.sites,p->target.provisioned,p->target.attached,p->target.active,p->target.sessions,p->target.rps,p->ramp,p->reach,p->hold);
        }
        printf("BFF operations: %zu; profile model: source-derived HTTP sessions; backend mutations: BFF only\n",catalog.operation_count);
        printf("Fixture provider: existing virtual node/factory/warehouse setup; UE capacity is not established by this plan.\n");
        goto done;
    }
    if(!strcmp(config.environment->status,"skip")) { printf("SKIP %s\n",config.name); goto done; }
    if(!wl_safe_path(opts->repo) || !wl_safe_path(opts->script_dir) || !wl_safe_path(opts->out_dir)) {
        rc=wl_error(&err,"workload repo/scripts/output paths must use letters, digits, / _ . : -"); goto done;
    }
    if(opts->run_id[0]) {
        if(!wl_safe_path(opts->run_id) || strchr(opts->run_id,'/')) { rc=wl_error(&err,"invalid run id");goto done; }
        strcpy(run_id,opts->run_id);
    } else snprintf(run_id,sizeof(run_id),"wl-%.100s-%ld-%ld",config.name,(long)time(NULL),(long)getpid());
    if(wl_path(run_dir,sizeof(run_dir),opts->out_dir,run_id,&err) || ulab_mkdir_p(opts->out_dir)) { rc=ULAB_EINTERNAL;goto done; }
    if(mkdir(run_dir,0700)) { rc=wl_error(&err,"cannot create %s: %s (workloads never overwrite an existing run)",run_dir,strerror(errno));goto done; }
    metrics=wl_metrics_open(run_dir,&config,&err); if(!metrics) {rc=ULAB_EINTERNAL;goto done;}
    char catalog_path[ULAB_MAX_PATH];
    if(wl_path(catalog_path,sizeof(catalog_path),run_dir,"catalog.snapshot.json",&err) || json_dump_file(catalog.manifest,catalog_path,JSON_INDENT(2))) {rc=wl_error(&err,"cannot snapshot operation catalog");goto done;}
    char profile_path[ULAB_MAX_PATH];json_t *profile=NULL;
    if(wl_path(profile_path,sizeof(profile_path),config.assets,"profiles/console.yaml",&err) || !(profile=wl_yaml_load(profile_path,&err))) {rc=ULAB_ERR;goto done;}
    if(wl_path(profile_path,sizeof(profile_path),run_dir,"profiles.snapshot.json",&err) || json_dump_file(profile,profile_path,JSON_INDENT(2))) rc=wl_error(&err,"cannot snapshot console profiles");
    json_decref(profile);if(rc)goto done;
    memset(&action,0,sizeof(action));action.sa_handler=stop_handler;sigemptyset(&action.sa_mask);
    stopping=0;sigaction(SIGINT,&action,&old_int);sigaction(SIGTERM,&action,&old_term);signals=1;
    if(wl_environment_open(&env,&config,opts,metrics,run_dir,&stopping,&err)) rc=ULAB_ERR;
    else rc=wl_schedule(&env,&catalog,&err);
    if(!rc && wl_environment_verify(&env,&err)) rc=ULAB_ERR;
    if(rc) stopping=1;
    ulab_copy(reason,sizeof(reason),err.msg);
    ulab_status("CLEANUP","%s",config.keep_environment?"retain workload environment":"release workload environment");
    cleanup_rc=wl_environment_cleanup(&env,&cleanup_err);
    ulab_progress_clear();
    if(cleanup_rc) { fprintf(stderr,"%s\n",cleanup_err.msg);if(!reason[0]) ulab_copy(reason,sizeof(reason),cleanup_err.msg); }
    wl_targets_t actual={0},target=config.phases[config.phase_count-1].target;
    if(env.site_ready && env.provisioned) wl_environment_counts(&env,&actual);
    wl_metrics_tick(metrics,"cleanup",&target,&actual,0,0,wl_now());
    rc=wl_metrics_finish(metrics,&config,rc,cleanup_rc,reason,config.assets,&err);
    reported=1;
done:
    ulab_progress_clear();
    if(err.msg[0] && !reported) fprintf(stderr,"%s\n",err.msg);
    if(signals) {sigaction(SIGINT,&old_int,NULL);sigaction(SIGTERM,&old_term,NULL);}
    wl_environment_close(&env);wl_metrics_close(metrics);wl_catalog_free(&catalog);wl_config_free(&config);
    return rc;
}
