# Ukama Engineering Principles 
*_listed in no particular order_ 

## Know The Business
You have a right to know how business functions and you have to understand it in order to make reasonable decisions. 

## Be open source at the core
Ukama is open source at its core. Our proprietary software must always use and may extend on the open-source version but never replace it. 

## [Be Cloud Native](https://github.com/cncf/toc/blob/main/DEFINITION.md#cncf-cloud-native-definition-v10)
Ukama modules and services should run in vendor-neutral, scalable, and immutable infrastructure. In other words, every module is a containerized microservice managed by an orchestrator. 

## Keep it simple ([KISS](https://en.wikipedia.org/wiki/KISS_principle))
System Design should be as simple as possible. Don’t overcomplicate without necessity. Simple systems are easier to understand, extend, and maintain.

## Do One Thing And Do It Well ([DOTADIW](https://en.wikipedia.org/wiki/Unix_philosophy))
Write programs that do one thing and do it really well. Write programs to work together. 

## Don’t add functionality without necessity ([YaGNI principle](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it))
In a dynamic environment, it’s impossible to predict requirements upfront so it does not make sense to spend time on something that is not needed right away

## [Prefer convention over Configuration](https://en.wikipedia.org/wiki/Convention_over_configuration)
A service or a system should run out of the box with the default configuration. 

## Follow Api First Approach
API is a first-class citizen. Internal user-facing applications should use the same API that is exposed externally. 

## [Prefer consistency in code and architecture](https://skiplist.com/insights/blog-software-principle-6-consistency-is-king)
When you face a choice like what naming convention or framework to use, try to look around, find similar cases and stay consistent with them unless you have a strong reason to diverge.

## Automate tests and infrastructure from the beginning 
Everything should be the code. Investments in automation early will pay off in the long run by simplifying the routine tasks and eliminating human error

## Security From Day 1
Think about security early from the beginning. It's hard to add security to an already existing system.

## [Fail early and visibly](https://en.wikipedia.org/wiki/Fail-fast)
In order to handle errors effectively, we need to know where they happen and what exactly happened. Systems are complex and they fail from time to time. Our goal is to make sure we are prepared 
