#include <stdio.h>
#include <string.h>
#include <jansson.h>

/* format: {[{"1", "/path/to/url"}, 
 *           {"3", "/other/path/"}]}
 */

#define AGENT_CB "agent-cbURL"
#define AGENT_TYPE "type"
#define AGENT_URL "url"

#define OCI "OCI-image"
#define PATH "http://registry.ukama.com/containers/etc/"

int main(void) {
  
  char* s = NULL;
  int count=3;

  json_t *root = json_object();
  json_object_set_new(root, AGENT_CB,  json_array());

  json_t *agentArray = json_object_get(root, AGENT_CB); 

  if (agentArray) {
    for (int i=0; i<count; i++) {
      json_t *agent = json_object();
      json_object_set_new(agent, AGENT_TYPE, json_string(OCI));
      json_object_set_new(agent, AGENT_URL, json_string(PATH));

      json_array_append(agentArray, agent);
    }
  }

  s = json_dumps(root, 0);
  
  puts(s);

  /* de-serial. */
  json_t *array = json_object_get(root, AGENT_CB);

  if (json_is_array(array)) {
    int size = json_array_size(array);

    for (int i=0; i<size; i++) {
      json_t *elem = json_array_get(array, i);
      json_t *type = json_object_get(elem, AGENT_TYPE);
      json_t *url  = json_object_get(elem, AGENT_URL);

      if (type && url) {
	printf("%d: type: %s url: %s\n", i, json_string_value(type),
	       json_string_value(url));
      }
    }
  }
    
  return 0;
}
