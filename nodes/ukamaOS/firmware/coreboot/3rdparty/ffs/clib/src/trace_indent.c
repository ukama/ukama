/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/src/trace_indent.c $                                     */
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

#include "trace_indent.h"
#include "stdio.h"

/*--------------------------------------------------
 * Private Constants
 --------------------------------------------------*/

#define DEFAULT_INDENT_SIZE   5
#define MAX_INDENT_SIZE       129
#define DEFAULT_INDENT        "     "

/*--------------------------------------------------
 * Global variables.
 --------------------------------------------------*/

/* The current trace level. */
unsigned char BB_TRACE_LEVEL = 0;
unsigned char BB_PREVIOUS_TRACE_LEVEL = 0;

/*--------------------------------------------------
 * Private global variables.
 --------------------------------------------------*/

static unsigned int Current_Indent_Level = 0;
static unsigned int Current_Indent_Size = DEFAULT_INDENT_SIZE;
static char Indent[MAX_INDENT_SIZE] = DEFAULT_INDENT;

/*--------------------------------------------------
 * Public Functions
 --------------------------------------------------*/

/*********************************************************************/
/*                                                                   */
/*   Function Name: Set_Indent_Size                                  */
/*                                                                   */
/*   Descriptive Name: Sets the number of spaces to use per indent.  */
/*                     If the current indent level is 5 and this     */
/*                     function is called to set the indent size to  */
/*                     10, then 50 spaces will be printed by the     */
/*                     trace macros before each line of trace        */
/*                     output at this indent level.                  */
/*                                                                   */
/*   Input: unsigned int Spaces_Per_Indent - A value less than 128.  */
/*                                                                   */
/*   Output: If Success : The specified value will be used for the   */
/*                        size of an indent.                         */
/*                                                                   */
/*           If Failure : This should only happen if the value       */
/*                        specified is >= 128, in which case it is   */
/*                        ignored and the previous value is retained.*/
/*                                                                   */
/*   Error Handling: Bad values for Spaces_Per_Indent are ignored.   */
/*                                                                   */
/*   Side Effects: The number of spaces per indent may be changed.   */
/*                                                                   */
/*   Notes:                                                          */
/*                                                                   */
/*********************************************************************/
void Set_Indent_Size(unsigned int Spaces_Per_Indent)
{
	unsigned int i;

	if (Spaces_Per_Indent <= MAX_INDENT_SIZE) {
		Current_Indent_Size = Spaces_Per_Indent;
		for (i = 0; i < Current_Indent_Size; i++) {
			Indent[i] = ' ';
		}
		Indent[i] = 0x00;
	}

	return;
}

/*********************************************************************/
/*                                                                   */
/*   Function Name: Indent_Trace_Output                              */
/*                                                                   */
/*   Descriptive Name: This function increases the current indent    */
/*                     level by one.                                 */
/*                                                                   */
/*   Input: None.                                                    */
/*                                                                   */
/*   Output: If Success : None                                       */
/*                                                                   */
/*           If Failure : None                                       */
/*                                                                   */
/*   Error Handling: None                                            */
/*                                                                   */
/*   Side Effects: The current indent level for trace output will    */
/*                 be increased by 1.                                */
/*                                                                   */
/*   Notes:                                                          */
/*                                                                   */
/*********************************************************************/
void Indent_Trace_Output(void)
{
	Current_Indent_Level += 1;
	return;
}

/*********************************************************************/
/*                                                                   */
/*   Function Name: Outdent_Trace_Output                             */
/*                                                                   */
/*   Descriptive Name: This function reduces the current indent      */
/*                     level if it is greater than 0.                */
/*                                                                   */
/*   Input: None                                                     */
/*                                                                   */
/*   Output: If Success : None                                       */
/*                                                                   */
/*           If Failure : None                                       */
/*                                                                   */
/*   Error Handling: If the current indent level is zero, then this  */
/*                   function does nothing.                          */
/*                                                                   */
/*   Side Effects: The current indent level for trace output may be  */
/*                 decreased by 1.                                   */
/*                                                                   */
/*   Notes:                                                          */
/*                                                                   */
/*********************************************************************/
void Outdent_Trace_Output(void)
{
	if (Current_Indent_Level > 0)
		Current_Indent_Level -= 1;
	return;
}

/*********************************************************************/
/*                                                                   */
/*   Function Name: Do_Indent                                        */
/*                                                                   */
/*   Descriptive Name: This function prints to stderr the number of  */
/*                     spaces corresponding to the current indent    */
/*                     level.                                        */
/*                                                                   */
/*   Input: None                                                     */
/*                                                                   */
/*   Output: If Success : None                                       */
/*                                                                   */
/*           If Failure : None                                       */
/*                                                                   */
/*   Error Handling: None                                            */
/*                                                                   */
/*   Side Effects: Some number of space may be output to stderr.     */
/*                                                                   */
/*   Notes:                                                          */
/*                                                                   */
/*********************************************************************/
void Do_Indent(void)
{
	unsigned int i;

	for (i = 0; i < Current_Indent_Level; i++) {
		fprintf(stderr, "%s", Indent);
	}

}
