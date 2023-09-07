/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/bb_trace.h $                                             */
/*                                                                        */
/* OpenPOWER FFS Project                                                  */
/*                                                                        */
/* Contributors Listed Below - COPYRIGHT 2014,2015                        */
/* [+] International Business Machines Corp.                              */
/*                                                                        */
/*                                                                        */
/* Licensed under the Apache License, Version 2.0 (the "License");        */
/* you may not use this file except in compliance with the License.       */
/* You may obtain a copy of the License at                                */
/*                                                                        */
/*     http://www.apache.org/licenses/LICENSE-2.0                         */
/*                                                                        */
/* Unless required by applicable law or agreed to in writing, software    */
/* distributed under the License is distributed on an "AS IS" BASIS,      */
/* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or        */
/* implied. See the License for the specific language governing           */
/* permissions and limitations under the License.                         */
/*                                                                        */
/* IBM_PROLOG_END_TAG                                                     */

#ifndef BB_TRACE_H

#define BB_TRACE_H

#include <stdio.h>
#include "trace_indent.h"

#if (ENABLE_API_TRACE || ENABLE_FUNCTION_TRACE || ENABLE_FUNCTION_PARAMETER_TRACE || ENABLE_API_TRACE || ENABLE_API_PARAMETER_TRACE || ENABLE_EXTRA_CHECKS || ENABLE_USER_TRACE )

#define ENABLE_BB_TRACE_BASE

#endif

#ifdef ENABLE_ALL_BB_TRACE

#define ENABLE_API_TRACE
#define ENABLE_FUNCTION_TRACE
#define ENABLE_API_PARAMETER_TRACE
#define ENABLE_FUNCTION_PARAMETER_TRACE
#define ENABLE_EXTRA_CHECKS
#define ENABLE_USER_TRACE
#define ENABLE_BB_TRACE_BASE

#endif

#ifdef ENABLE_BB_TRACE_BASE

  /* The trace level is a bitmask.  The following bits, if 1, activate the corresponding trace function:
   *
   *  0 : API Entry & Exit
   *  1 : Local function entry & exit
   *  2 : API parameters will be listed on entry and return values listed on exit
   *  3 : Local function parameters and return values will be listed
   *  4 : Low bit, extra checks level number.
   *  5 : High bit, extra checks level number.
   *  6 : Low bit, user defined trace level number.
   *  7 : High bit, user defined trace level number.
   */

#define API_ENTRY_EXIT_TRACE              0x01
#define FUNCTION_ENTRY_EXIT_TRACE         0x02
#define API_PARAMETER_TRACE               0x04
#define FUNCTION_PARAMETER_TRACE          0x08
#define EXTRA_CHECKS_LEVEL_1              0x10
#define EXTRA_CHECKS_LEVEL_2              0x20
#define EXTRA_CHECKS_LEVEL_3              0x30
#define USER_TRACE_LEVEL_1                0x40
#define USER_TRACE_LEVEL_2                0x80
#define USER_TRACE_LEVEL_3                0xC0

extern unsigned char BB_TRACE_LEVEL;
extern unsigned char BB_PREVIOUS_TRACE_LEVEL;

#define TRACE_INDENT()                   if (BB_TRACE_LEVEL != 0)                                                      \
                                             Indent_Trace_Output();

#define TRACE_OUTDENT()                  if ( BB_TRACE_LEVEL != 0 )                                                    \
                                             Outdent_Trace_Output();

#define TRACE_STOP()                     if ( BB_TRACE_LEVEL != 0 )                                                    \
                                           {                                                                             \
                                             BB_PREVIOUS_TRACE_LEVEL = BB_TRACE_LEVEL;                                   \
                                             BB_TRACE_LEVEL = 0;                                                         \
                                           }

#define TRACE_START()                    if ( BB_TRACE_LEVEL == 0 )                                                    \
                                           {                                                                             \
                                             BB_TRACE_LEVEL = BB_PREVIOUS_TRACE_LEVEL;                                   \
                                           }

#define SET_TRACE_LEVEL( Level )         BB_TRACE_LEVEL = Level;

#else

#define API_ENTRY_EXIT_TRACE
#define FUNCTION_ENTRY_EXIT_TRACE
#define API_PARAMETER_TRACE
#define FUNCTION_PARAMETER_TRACE
#define EXTRA_CHECKS_LEVEL_1
#define EXTRA_CHECKS_LEVEL_2
#define EXTRA_CHECKS_LEVEL_3
#define USER_TRACE_LEVEL_1
#define USER_TRACE_LEVEL_2
#define USER_TRACE_LEVEL_3

#define TRACE_INDENT()

#define TRACE_OUTDENT()

#define TRACE_STOP()

#define TRACE_START()

#define SET_TRACE_LEVEL( Level )

#endif

#ifdef ENABLE_FUNCTION_TRACE

#define FUNCTION_ENTRY(  )                     if ( BB_TRACE_LEVEL & FUNCTION_ENTRY_EXIT_TRACE )                       \
                                                 {                                                                       \
                                                    fprintf( stderr, " \n\n");                                           \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "-------------------- \n");                         \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "Entering %s.\n", __func__ );                       \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "-------------------- \n\n");                       \
                                                    Indent_Trace_Output();                                               \
                                                 }

#define FUNCTION_EXIT(  )                      if ( BB_TRACE_LEVEL & FUNCTION_ENTRY_EXIT_TRACE )                       \
                                                 {                                                                       \
                                                    fprintf( stderr, " \n\n");                                           \
                                                    Outdent_Trace_Output();                                              \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "-------------------- \n");                         \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "Leaving %s.\n", __func__ );                        \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "-------------------- \n\n");                       \
                                                 }                                                                       \
                                                 return;

#define FUNCTION_RETURN( Format, Value )                                                                                             \
                                                 if ( BB_TRACE_LEVEL & FUNCTION_ENTRY_EXIT_TRACE )                                     \
                                                 {                                                                                     \
                                                    fprintf( stderr, " \n\n");                                                         \
                                                    Outdent_Trace_Output();                                                            \
                                                    Do_Indent();                                                                       \
                                                    fprintf( stderr, "-------------------- \n");                                       \
                                                    Do_Indent();                                                                       \
                                                    fprintf( stderr, "Leaving %s with return value " Format "\n", __func__ , Value);   \
                                                    Do_Indent();                                                                       \
                                                    fprintf( stderr, "-------------------- \n\n");                                     \
                                                 }                                                                                     \
                                                 return Value;

#else

#define FUNCTION_ENTRY( )

#define FUNCTION_EXIT( )                       return;

#define FUNCTION_RETURN( Format, Value )       return Value;

#endif

#ifdef ENABLE_FUNCTION_PARAMETER_TRACE

#define PRINT_FUNCTION_PARAMETER( ... )        if ( BB_TRACE_LEVEL & FUNCTION_PARAMETER_TRACE )                        \
                                                 {                                                                       \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr,  __VA_ARGS__ );                                      \
                                                   fprintf( stderr, "\n");                                               \
                                                 }

#else

#define PRINT_FUNCTION_PARAMETER( ... )

#endif

#ifdef ENABLE_API_TRACE

#define API_FUNCTION_ENTRY( )                  if ( BB_TRACE_LEVEL & API_ENTRY_EXIT_TRACE )                            \
                                                 {                                                                       \
                                                   fprintf( stderr, " \n\n");                                            \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "==================== \n");                         \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "Entering %s API.\n", __func__ );                   \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "==================== \n\n");                       \
                                                    Indent_Trace_Output();                                               \
                                                 }

#define API_FUNCTION_EXIT( )                   if ( BB_TRACE_LEVEL & API_ENTRY_EXIT_TRACE )                            \
                                                 {                                                                       \
                                                   fprintf( stderr, " \n\n");                                            \
                                                    Outdent_Trace_Output();                                              \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "==================== \n");                         \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "Leaving %s API.\n", __func__ );                    \
                                                    Do_Indent();                                                         \
                                                    fprintf( stderr, "==================== \n\n");                       \
                                                 }                                                                       \
                                                 return;

#define API_FUNCTION_RETURN( Format, Value )                                                                                            \
                                                 if ( BB_TRACE_LEVEL & API_ENTRY_EXIT_TRACE )                                             \
                                                 {                                                                                        \
                                                   fprintf( stderr, " \n\n");                                                             \
                                                    Outdent_Trace_Output();                                                               \
                                                    Do_Indent();                                                                          \
                                                    fprintf( stderr, "==================== \n");                                          \
                                                    Do_Indent();                                                                          \
                                                    fprintf( stderr, "Leaving %s API with return value " Format "\n", __func__, Value);   \
                                                    Do_Indent();                                                                          \
                                                    fprintf( stderr, "==================== \n\n");                                        \
                                                 }                                                                                        \
                                                 return Value;

#else

#define API_FUNCTION_ENTRY( )

#define API_FUNCTION_EXIT( )                   return;

#define API_FUNCTION_RETURN( Format, Value )   return Value;

#endif

#ifdef ENABLE_API_PARAMETER_TRACE

#define PRINT_API_PARAMETER( ... )             if ( BB_TRACE_LEVEL & API_PARAMETER_TRACE )                             \
                                                 {                                                                       \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr,  __VA_ARGS__ );                                      \
                                                   fprintf( stderr, "\n");                                               \
                                                 }

#else

#define PRINT_API_PARAMETER( ... )

#endif

#ifdef  ENABLE_EXTRA_CHECKS

#define LEVEL1_EXTRA_CHECK( ... )         if ( BB_TRACE_LEVEL & EXTRA_CHECKS_LEVEL_1 )                                 \
                                            {                                                                            \
                                              if ( __VA_ARGS__ )                                                         \
                                              {                                                                          \
                                                   fprintf( stderr, "\n");                                               \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "Extra check in file %s, line %d has failed!\n",     \
                                                            __FILE__, __LINE__ );                                        \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                              }                                                                          \
                                                                                                                         \
                                            }

#define LEVEL1_EC_PRINT_LINE( ... )       if ( BB_TRACE_LEVEL & EXTRA_CHECKS_LEVEL_1 )                                 \
                                            {                                                                            \
                                               Do_Indent();                                                              \
                                               fprintf( stderr,  __VA_ARGS__ );                                          \
                                               fprintf( stderr, "\n");                                                   \
                                            }

#define LEVEL2_EXTRA_CHECK( ... )         if ( BB_TRACE_LEVEL & EXTRA_CHECKS_LEVEL_2 )                                 \
                                            {                                                                            \
                                              if ( __VA_ARGS__ )                                                         \
                                              {                                                                          \
                                                   fprintf( stderr, "\n");                                               \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "Extra check in file %s, line %d has failed!\n",     \
                                                            __FILE__, __LINE__ );                                        \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                              }                                                                          \
                                                                                                                         \
                                            }

#define LEVEL2_EC_PRINT_LINE( ... )       if ( BB_TRACE_LEVEL & EXTRA_CHECKS_LEVEL_2 )                                 \
                                            {                                                                            \
                                               Do_Indent();                                                              \
                                               fprintf( stderr,  __VA_ARGS__ );                                          \
                                               fprintf( stderr, "\n");                                                   \
                                            }

#define LEVEL3_EXTRA_CHECK( ... )         if ( BB_TRACE_LEVEL & EXTRA_CHECKS_LEVEL_3 )                                 \
                                            {                                                                            \
                                              if ( __VA_ARGS__ )                                                         \
                                              {                                                                          \
                                                   fprintf( stderr, "\n");                                               \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "Extra check in file %s, line %d has failed!\n",     \
                                                            __FILE__, __LINE__ );                                        \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                              }                                                                          \
                                                                                                                         \
                                            }

#define LEVEL3_EC_PRINT_LINE( ... )       if ( BB_TRACE_LEVEL & EXTRA_CHECKS_LEVEL_3 )                                 \
                                            {                                                                            \
                                               Do_Indent();                                                              \
                                               fprintf( stderr,  __VA_ARGS__ );                                          \
                                               fprintf( stderr, "\n");                                                   \
                                            }

#define EXTRA_CHECK( Trace_Level, ... )   if ( BB_TRACE_LEVEL & Trace_Level )                                          \
                                            {                                                                            \
                                              if ( __VA_ARGS__ )                                                         \
                                              {                                                                          \
                                                   fprintf( stderr, "\n");                                               \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "Extra check in file %s, line %d has failed!\n",     \
                                                            __FILE__, __LINE__ );                                        \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr, "!!!!!!!!!!!!!!!!!!!! \n");                          \
                                              }                                                                          \
                                                                                                                         \
                                            }

#else

#define LEVEL1_EXTRA_CHECK( ... )

#define LEVEL2_EXTRA_CHECK( ... )

#define LEVEL3_EXTRA_CHECK( ... )

#define LEVEL1_EC_PRINT_LINE( ... )

#define LEVEL2_EC_PRINT_LINE( ... )

#define LEVEL3_EC_PRINT_LINE( ... )

#define EXTRA_CHECK( Trace_Level, ... )

#endif

#ifdef  ENABLE_USER_TRACE

#define USER1_PRINT_LINE( ... )                if ( BB_TRACE_LEVEL & USER_TRACE_LEVEL_1 )                              \
                                                 {                                                                       \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr,  __VA_ARGS__ );                                      \
                                                   fprintf( stderr, "\n");                                               \
                                                 }

#define USER2_PRINT_LINE( ... )                if ( BB_TRACE_LEVEL & USER_TRACE_LEVEL_2 )                              \
                                                 {                                                                       \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr,  __VA_ARGS__ );                                      \
                                                   fprintf( stderr, "\n");                                               \
                                                 }

#define USER3_PRINT_LINE( ... )                if ( BB_TRACE_LEVEL & USER_TRACE_LEVEL_3 )                              \
                                                 {                                                                       \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr,  __VA_ARGS__ );                                      \
                                                   fprintf( stderr, "\n");                                               \
                                                 }

#define TRACE_PRINT_LINE( Trace_Level, ... )   if ( BB_TRACE_LEVEL & Trace_Level )                                     \
                                                 {                                                                       \
                                                   Do_Indent();                                                          \
                                                   fprintf( stderr,  __VA_ARGS__ );                                      \
                                                   fprintf( stderr, "\n");                                               \
                                                 }

#else

#define USER1_PRINT_LINE( ... )

#define USER2_PRINT_LINE( ... )

#define USER3_PRINT_LINE( ... )

#define TRACE_PRINT_LINE( Trace_Level, ... )

#endif

#endif
