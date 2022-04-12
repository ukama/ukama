--
-- Copyright (C) 2015-2016 secunet Security Networks AG
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

with HW.GFX.EDID;

package HW.GFX.GMA.Display_Probing
is

   type Port_List_Range is range 0 .. 7;
   type Port_List is array (Port_List_Range) of Port_Type;
   All_Ports : constant Port_List :=
     (DP1, DP2, DP3, HDMI1, HDMI2, HDMI3, Analog, Internal);

   procedure Read_EDID
     (Raw_EDID :    out EDID.Raw_EDID_Data;
      Port     : in     Active_Port_Type;
      Success  :    out Boolean)
   with
      Post => (if Success then EDID.Valid (Raw_EDID));

   procedure Scan_Ports
     (Configs     :    out Pipe_Configs;
      Ports       : in     Port_List := All_Ports;
      Max_Pipe    : in     Pipe_Index := Pipe_Index'Last;
      Keep_Power  : in     Boolean := False);

   procedure Hotplug_Events (Ports : out Port_List);

end HW.GFX.GMA.Display_Probing;
