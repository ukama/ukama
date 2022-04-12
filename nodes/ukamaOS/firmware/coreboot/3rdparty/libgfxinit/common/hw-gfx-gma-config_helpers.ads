--
-- Copyright (C) 2015-2017 secunet Security Networks AG
--
-- This program is free software; you can redistribute it and/or modify
-- it under the terms of the GNU General Public License as published by
-- the Free Software Foundation; either version 2 of the License, or
-- (at your option) any later version.
--
-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU General Public License for more details.
--

with HW;

private package HW.GFX.GMA.Config_Helpers
is

   function To_GPU_Port
     (Pipe  : Pipe_Index;
      Port  : Active_Port_Type)
      return GPU_Port;

   function To_PCH_Port (Port : Active_Port_Type) return PCH_Port;

   function To_Display_Type (Port : Active_Port_Type) return Display_Type;

   procedure Fill_Port_Config
     (Port_Cfg :    out Port_Config;
      Pipe     : in     Pipe_Index;
      Port     : in     Port_Type;
      Mode     : in     Mode_Type;
      Success  :    out Boolean)
   with
      Post =>
        (if Success then
            Port_Cfg.Mode.H_Visible = Mode.H_Visible and
            Port_Cfg.Mode.V_Visible = Mode.V_Visible);

   ----------------------------------------------------------------------------

   pragma Warnings (GNAT, Off, """Integer_32"" is already use-visible *",
                    Reason => "Needed for older compiler versions");
   use type HW.Pos32;
   pragma Warnings (GNAT, On, """Integer_32"" is already use-visible *");

   -- Validate just enough to satisfy Pipe_Setup pre conditions.
   function Valid_FB
     (FB    : Framebuffer_Type;
      Mode  : Mode_Type)
      return Boolean is
     (Rotated_Width (FB) <= Mode.H_Visible and
      Rotated_Height (FB) <= Mode.V_Visible and
      (FB.Offset = VGA_PLANE_FRAMEBUFFER_OFFSET or
       FB.Height + FB.Start_Y <= FB.V_Stride));

   -- Also validate that we only use supported values / features.
   function Validate_Config
     (FB                : Framebuffer_Type;
      Mode              : Mode_Type;
      Pipe              : Pipe_Index)
      return Boolean
   with
      Post => (if Validate_Config'Result then Valid_FB (FB, Mode));

   -- For still active pipes, ensure only timings
   -- changed that don't affect FB validity.
   function Stable_FB (Old_C, New_C : Pipe_Configs) return Boolean is
     (for all P in Pipe_Index =>
         New_C (P).Port = Disabled or
           (New_C (P).Port = Old_C (P).Port and
            New_C (P).Framebuffer = Old_C (P).Framebuffer and
            New_C (P).Cursor = Old_C (P).Cursor and
            New_C (P).Mode.H_Visible = Old_C (P).Mode.H_Visible and
            New_C (P).Mode.V_Visible = Old_C (P).Mode.V_Visible));

   ----------------------------------------------------------------------------

   function Highest_Dotclock (Configs : Pipe_Configs) return Frequency_Type;

   procedure Limit_Dotclocks
     (Configs  : in out Pipe_Configs;
      Max      : in     Frequency_Type)
   with
      Post => Stable_FB (Configs'Old, Configs);

end HW.GFX.GMA.Config_Helpers;
