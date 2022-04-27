import {message, danger} from "danger"

const pr = danger.github.pr


// No PR is too small to include a description of why you made a change
if (pr.mergeable_state != 'draft' ){
  if (pr.body === null || pr.body.length < 20) {
    fail('Please include a description of your PR changes.');
  } else {
    var match = pr.body.match(/\bResolves #(\d*)\b/i);
    if (match === null){
      match = pr.body.match(/\Fixes #(\d*)\b/i);
    }
    if (match !== null){
        const issueNumber = match.pop();
        const intIssueNum = parseInt(issueNumber)
        if (issueNumber === undefined || intIssueNum === undefined){
            warn("Error parsing the issue number");
            return;
        }
      
        danger.github.api.issues.get({
            issue_number: intIssueNum, 
            owner: danger.github.thisPR.owner,
            repo: danger.github.thisPR.repo,
          }).catch(error => {
          if (error.status === 404) {
            fail(`Cannot find issue #${issueNumber}`);
          } else if (error.status !== 200) {
            warn(`Unable to check issue #${intIssueNum}. Error status: ${error.status}. Error data: ${error}`);
          }
          })
        
    }else{
        fail("PR does not have a reference to an issue in the description. Consider adding something like `Fixes #123`");
    }
  }


// // Always ensure we assign someone to a PR, if its a
// if (pr.assignee === null) {
//     //const method = pr.title.includes("WIP") ? warn : fail
//     warn("Please assign someone to this PR.");
// }



}