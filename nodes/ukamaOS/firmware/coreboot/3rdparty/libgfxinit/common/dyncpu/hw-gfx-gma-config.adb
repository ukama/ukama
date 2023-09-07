--
-- Copyright (C) 2018 Nico Huber <nico.h@gmx.de>
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

package body HW.GFX.GMA.Config
is

   procedure Detect_CPU (Device : Word16) is
   begin
      for CPU in Gen_CPU_Type loop
         for CPU_Var in Gen_CPU_Variant loop
            if Is_GPU (Device, CPU, CPU_Var) then
               Config.CPU := CPU;
               Config.CPU_Var := CPU_Var;
               exit;
            end if;
         end loop;
      end loop;
   end Detect_CPU;

end HW.GFX.GMA.Config;
