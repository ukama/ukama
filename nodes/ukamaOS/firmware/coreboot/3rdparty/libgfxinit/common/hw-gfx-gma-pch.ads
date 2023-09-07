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

with HW.GFX.GMA.Config;

private package HW.GFX.GMA.PCH is

   type FDI_Port_Type is (FDI_A, FDI_B, FDI_C);

   ----------------------------------------------------------------------------

   -- common to all PCH outputs

   function PCH_TRANSCODER_SELECT_SHIFT return Natural is
     (if Config.Has_New_FDI_Sink then 29 else 30);

   function PCH_TRANSCODER_SELECT_MASK return Word32 is
     (if Config.Has_New_FDI_Sink then 3 * 2 ** 29 else 1 * 2 ** 30);

   function PCH_TRANSCODER_SELECT (Port : FDI_Port_Type) return Word32 is
     (case Port is
         when FDI_A => Shift_Left (0, PCH_TRANSCODER_SELECT_SHIFT),
         when FDI_B => Shift_Left (1, PCH_TRANSCODER_SELECT_SHIFT),
         when FDI_C => Shift_Left (2, PCH_TRANSCODER_SELECT_SHIFT));

end HW.GFX.GMA.PCH;
