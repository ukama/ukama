/* SPDX-License-Identifier: MPL-2.0 */
#include "workload.h"
#include <ctype.h>
#include <string.h>

static const char *string(const json_t *o,const char *key) {
    const char *s=json_string_value(json_object_get(o,key));return s?s:"";
}
static double number(const json_t *o,const char *key) { return json_number_value(json_object_get(o,key)); }
static unsigned long long count(const json_t *o,const char *key) { return (unsigned long long)json_integer_value(json_object_get(o,key)); }
/* Error details can contain upstream newlines/control characters. Keep each
 * report entry on one line, without letting a response move the terminal. */
static void plain(FILE *out,const char *s) {
    for(const unsigned char *p=(const unsigned char *)s;*p;p++) fputc(iscntrl(*p)?' ':*p,out);
}
static size_t limit(size_t n,size_t cap) { return cap && n>cap?cap:n; }
static void stage_write(FILE *out,const json_t *stage,const json_t *thresholds) {
    const char *status=string(stage,"status");json_t *checks=json_object_get(stage,"checks");
    fprintf(out,"%-10s %s: operation_checks=%zu/%zu duration=%.1fs\n",status,string(stage,"phase"),
        json_array_size(checks)-(size_t)count(stage,"failed_checks"),json_array_size(checks),number(stage,"duration_sec"));
    if(json_array_size(checks)) {
        fprintf(out,"           errors expected<=%g%% actual=%.2f%%; missed expected<=%g%% actual=%.2f%%; p95 ",
            number(thresholds,"max_error_percent"),number(stage,"worst_error_percent"),
            number(thresholds,"max_missed_percent"),number(stage,"worst_missed_percent"));
        if(number(thresholds,"p95_ms")>0) fprintf(out,"expected<=%gms ",number(thresholds,"p95_ms"));else fputs("limit=disabled ",out);
        if(count(stage,"completed")) fprintf(out,"actual=%.2fms\n",number(stage,"worst_p95_ms"));else fputs("actual=n/a\n",out);
        fprintf(out,"           reads=%llu ok=%llu failed=%llu missed=%llu throttled=%llu timeouts=%llu successful_rate=%.2f/s\n",
            count(stage,"completed"),count(stage,"successful"),count(stage,"failed"),count(stage,"missed"),
            count(stage,"throttled"),count(stage,"timeouts"),number(stage,"successful_rps"));
    }
    if(strcmp(status,"PASS") || !json_array_size(checks)) {
        fputs("           ",out);plain(out,string(stage,"reason"));fputc('\n',out);
    }
    if(json_is_true(json_object_get(stage,"entered"))) {
        json_t *expected=json_object_get(stage,"expected_population"),*actual=json_object_get(stage,"actual_population");
        fprintf(out,"           population expected/actual_at_end sites=%llu/%llu sims=%llu/%llu attached=%llu/%llu active=%llu/%llu sessions=%llu/%llu\n",
            count(expected,"sites"),count(actual,"sites"),count(expected,"provisioned_ues"),count(actual,"provisioned_ues"),
            count(expected,"attached_ues"),count(actual,"attached_ues"),count(expected,"active_ues"),count(actual,"active_ues"),
            count(expected,"console_sessions"),count(actual,"console_sessions"));
    }
}
static void phases_write(FILE *out,const json_t *s,const json_t *thresholds) {
    json_t *phases=json_object_get(s,"phase_results"),*counts=json_object_get(s,"phase_counts");
    if(!phases) return;
    fprintf(out,"\nPHASES     passed=%llu failed=%llu incomplete=%llu not_run=%llu no_data=%llu (setup separate)\n",
        count(counts,"PASS"),count(counts,"FAIL"),count(counts,"INCOMPLETE"),count(counts,"NOT_RUN"),count(counts,"NO_DATA"));
    fprintf(out,"%-24s %9s %-11s %-11s %-11s\n","PHASE","TARGET/s","REACH","HOLD","RESULT");
    size_t i;json_t *p;
    json_array_foreach(phases,i,p)
        fprintf(out,"%-24s %9g %-11s %-11s %-11s\n",string(p,"name"),number(json_object_get(p,"target"),"rps"),
            string(json_object_get(p,"reach"),"status"),string(json_object_get(p,"hold"),"status"),string(p,"status"));
    fputs("\nPhase checks in scenario order (reach includes ramp/readiness; hold is steady load).\n",out);
    fputs("Actual percentages/p95 below are the WORST operation/category values, not blended averages.\n",out);
    json_t *setup=json_object_get(s,"setup_result");
    if(setup) {fprintf(out,"\nPHASE      setup result=%s\n",string(setup,"status"));stage_write(out,json_object_get(setup,"reach"),thresholds);}
    json_array_foreach(phases,i,p) {
        fprintf(out,"\nPHASE      %s result=%s rate=%g->%g/s ramp=%gs hold=%gs\n",string(p,"name"),string(p,"status"),
            number(json_object_get(p,"from"),"rps"),number(json_object_get(p,"target"),"rps"),number(p,"ramp_seconds"),number(p,"hold_seconds"));
        stage_write(out,json_object_get(p,"reach"),thresholds);stage_write(out,json_object_get(p,"hold"),thresholds);
    }
}

int wl_report_write(FILE *out,const json_t *s,size_t max_details) {
    json_t *thresholds=json_object_get(s,"thresholds"),*latency=json_object_get(s,"read_latency");
    json_t *ops=json_object_get(s,"operation_totals"),*violations=json_object_get(s,"threshold_violations");
    json_t *examples=json_object_get(s,"error_examples");
    const char *verdict=json_is_true(json_object_get(s,"passed"))?"PASS":"FAIL";
    fprintf(out,"\nREPORT     workload summary\nSCENARIO   %s\n",string(s,"scenario"));
    fprintf(out,"RESULT     %s duration=%.1fs\n",verdict,number(s,"duration_sec"));
    fprintf(out,"READS      completed=%llu successful=%llu failed=%llu (%.2f%%) missed=%llu (%.2f%%)\n",
        count(s,"requests"),count(s,"requests")-count(s,"failed_requests"),count(s,"failed_requests"),
        number(s,"error_percent"),count(s,"missed_requests"),number(s,"missed_percent"));
    if(count(s,"requests"))
        fprintf(out,"LATENCY    p50=%.2fms p95=%.2fms p99=%.2fms max=%.2fms (all completed reads)\n",
            number(latency,"p50_ms"),number(latency,"p95_ms"),number(latency,"p99_ms"),number(latency,"max_ms"));
    else fputs("LATENCY    no completed reads\n",out);
    fprintf(out,"CHECKS     execution=%s cleanup=%s metrics=%s threshold_violations=%zu\n",
        string(s,"execution"),string(s,"cleanup"),json_is_true(json_object_get(s,"metrics_complete"))?"complete":"incomplete",json_array_size(violations));
    fprintf(out,"LIMITS     per phase/operation/category: errors<=%g%% missed<=%g%% p95",
        number(thresholds,"max_error_percent"),number(thresholds,"max_missed_percent"));
    if(number(thresholds,"p95_ms")>0) fprintf(out,"<=%gms\n",number(thresholds,"p95_ms"));
    else fputs("=disabled\n",out);
    if(*string(s,"reason")) { fputs("REASON     ",out);plain(out,string(s,"reason"));fputc('\n',out); }
    phases_write(out,s,thresholds);
    if(json_array_size(ops)) {
        fputs("\nRead operations across all phases (percentiles merged from histograms):\n",out);
        fprintf(out,"%-24s %-7s %8s %6s %6s %7s %9s\n","OPERATION","ACTOR","DONE","FAIL","MISS","ERR%","P95 ms");
        size_t i;json_t *r;
        json_array_foreach(ops,i,r)
            fprintf(out,"%-24s %-7s %8llu %6llu %6llu %7.2f %9.2f\n",string(r,"operation"),string(r,"category"),
                count(r,"completed"),count(r,"failed"),count(r,"missed"),number(r,"error_percent"),number(r,"p95_ms"));
    }
    size_t n=json_array_size(violations);
    if(!json_array_size(ops)) fputs("\nTHRESHOLDS no console/rate samples; no threshold checks performed\n",out);
    else if(!n) fputs("\nTHRESHOLDS all observed operation/phase checks passed\n",out);
    for(size_t i=0;i<limit(n,max_details);i++) {
        json_t *v=json_array_get(violations,i);
        fprintf(out,"\nVIOLATION  %s %s/%s\n",string(v,"phase"),string(v,"category"),string(v,"operation"));
        if(number(v,"error_percent")>number(thresholds,"max_error_percent"))
            fprintf(out,"           errors %.2f%% > %g%% (%llu/%llu completed)\n",number(v,"error_percent"),number(thresholds,"max_error_percent"),count(v,"failed"),count(v,"completed"));
        if(number(v,"missed_percent")>number(thresholds,"max_missed_percent"))
            fprintf(out,"           missed %.2f%% > %g%% (%llu missed)\n",number(v,"missed_percent"),number(thresholds,"max_missed_percent"),count(v,"missed"));
        if(number(thresholds,"p95_ms")>0 && number(v,"p95_ms")>number(thresholds,"p95_ms"))
            fprintf(out,"           p95 %.2fms > %gms\n",number(v,"p95_ms"),number(thresholds,"p95_ms"));
    }
    if(n>limit(n,max_details)) fprintf(out,"           %zu more violations in report.txt / report.json / workload.html\n",n-limit(n,max_details));
    n=json_array_size(examples);
    if(n) fputs("\nFirst error per affected phase/operation/category (all attempts in requests.jsonl):\n",out);
    for(size_t i=0;i<limit(n,max_details);i++) {
        json_t *v=json_array_get(examples,i);
        fprintf(out,"ERROR      %s %s/%s outcome=%s http=%llu\n",string(v,"phase"),string(v,"category"),string(v,"operation"),string(v,"outcome"),count(v,"http_status"));
        fputs("           ",out);plain(out,string(v,"detail"));fputc('\n',out);
        fprintf(out,"           request_id=%s\n",string(v,"request_id"));
    }
    if(n>limit(n,max_details)) fprintf(out,"           %zu more examples in the report files\n",n-limit(n,max_details));
    json_t *writes=json_object_get(s,"report_write_errors");
    for(size_t i=0;i<json_array_size(writes);i++) {
        fputs("ARTIFACT   ",out);plain(out,json_string_value(json_array_get(writes,i)));fputc('\n',out);
    }
    fprintf(out,"\n%s       %s\nREPORT     %s/workload.html\nJSON       %s/report.json\nTEXT       %s/report.txt\n",
        verdict,string(s,"scenario"),string(s,"output_dir"),string(s,"output_dir"),string(s,"output_dir"));
    return ferror(out)?ULAB_ERR:ULAB_OK;
}
