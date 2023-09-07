--
-- Copyright (C) 2016 secunet Security Networks AG
-- Copyright (C) 2019 Nico Huber <nico.h@gmx.de>
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

with HW.Time;
with HW.GFX.GMA.Config;
with HW.GFX.GMA.Registers;

package body HW.GFX.GMA.Power_And_Clocks is

   FSB_FREQ_SEL_MASK : constant := 7 * 2 ** 0;
   CLKCFG_FSB_400    : constant Frequency_Type := 100_000_000;
   CLKCFG_FSB_533    : constant Frequency_Type := 133_333_333;
   CLKCFG_FSB_667    : constant Frequency_Type := 166_666_666;
   CLKCFG_FSB_800    : constant Frequency_Type := 200_000_000;
   CLKCFG_FSB_1067   : constant Frequency_Type := 266_666_666;
   CLKCFG_FSB_1333   : constant Frequency_Type := 333_333_333;

   type Div_Array is array (0 .. 7) of Pos64;

   procedure Get_VCO (VCO : out Int64; Divisors : out Div_Array)
   is
      G45_3200 : constant Div_Array := (12, 10,  8,  7,  5, 16, others => 1);
      G45_4000 : constant Div_Array := (14, 12, 10,  8,  6, 20, others => 1);
      G45_4800 : constant Div_Array := (20, 14, 12, 10,  8, 24, others => 1);
      G45_5333 : constant Div_Array := (20, 16, 12, 12,  8, 28, others => 1);
      G45_Divs : constant array (Natural range 0 .. 7) of Div_Array :=
        (G45_3200, G45_4000, G45_5333, G45_4800, others => (others => 1));

      GM45_2667 : constant Div_Array := (12,  8, others => 1);
      GM45_3200 : constant Div_Array := (14, 10, others => 1);
      GM45_4000 : constant Div_Array := (18, 12, others => 1);
      GM45_5333 : constant Div_Array := (24, 16, others => 1);
      GM45_Divs : constant array (Natural range 0 .. 7) of Div_Array :=
        (0 => GM45_3200, 1 => GM45_4000, 2 => GM45_5333, 4 => GM45_2667,
         others => (others => 1));

      HPLLVCO : Word32;
      VCO_Sel : Natural range 0 .. 7;
   begin
      if Config.Has_GMCH_Mobile_VCO then
         Registers.Read (Registers.GMCH_HPLLVCO_MOBILE, HPLLVCO);
         VCO_Sel := Natural (HPLLVCO and 7);
         VCO :=
           (case VCO_Sel is
               when 0 => 3_200_000_000,
               when 1 => 4_000_000_000,
               when 2 => 5_333_333_333,
               --when 3 => 6_400_000_000,
               when 4 => 2_666_666_667,
               --when 5 => 4_266_666_667,
               when others => 0);
         Divisors := GM45_Divs (VCO_Sel);
      else
         Registers.Read (Registers.GMCH_HPLLVCO, HPLLVCO);
         VCO_Sel := Natural (HPLLVCO and 7);
         VCO :=
           (case VCO_Sel is
               when 0 => 3_200_000_000,
               when 1 => 4_000_000_000,
               when 2 => 5_333_333_333,
               when 3 => 4_800_000_000,
               when others => 0);
         Divisors := G45_Divs (VCO_Sel);
      end if;
   end Get_VCO;

   procedure Get_CDClk (CDClk : out Config.CDClk_Range)
   is
      use type HW.Word16;

      Tmp_Clk : Int64 := 0;

      VCO : Int64;
      Divisors : Div_Array;

      GCFGC : Word16;
      CDClk_Sel : Natural range 0 .. 7;
   begin
      if PCI_Usable then
         Get_VCO (VCO, Divisors);
         PCI_Read16 (GCFGC, 16#f0#);
         if Config.Has_GMCH_Mobile_VCO then
            CDClk_Sel := Natural (Shift_Right (GCFGC, 12) and 1);
         else
            CDClk_Sel := Natural (Shift_Right (GCFGC, 4) and 7);
         end if;
         Tmp_Clk := VCO / Divisors (CDClk_Sel);
      end if;

      if Tmp_Clk in Config.CDClk_Range then
         CDClk := Tmp_Clk;
      else
         if Config.Has_GMCH_Mobile_VCO then
            CDClk := 5_333_333_333 / 24;
         else
            CDClk := 5_333_333_333 / 28;
         end if;
      end if;
   end Get_CDClk;

   -- The Raw Freq is 1/4 of the FSB freq
   procedure Get_Raw_Clock (Raw_Clock : out Frequency_Type)
   is
      CLK_CFG : Word32;
      type Freq_Sel is new Natural range 0 .. 7;
   begin
      Registers.Read
        (Register => Registers.GMCH_CLKCFG,
         Value => CLK_CFG);
      case Freq_Sel (CLK_CFG and FSB_FREQ_SEL_MASK) is
         when 0      => Raw_Clock := CLKCFG_FSB_1067;
         when 1      => Raw_Clock := CLKCFG_FSB_533;
         when 2      => Raw_Clock := CLKCFG_FSB_800;
         when 3      => Raw_Clock := CLKCFG_FSB_667;
         when 4      => Raw_Clock := CLKCFG_FSB_1333;
         when 5      => Raw_Clock := CLKCFG_FSB_400;
         when 6      => Raw_Clock := CLKCFG_FSB_1067;
         when 7      => Raw_Clock := CLKCFG_FSB_1333;
      end case;
   end Get_Raw_Clock;

   procedure Initialize
   is
      CDClk : Config.CDClk_Range;
   begin
      Get_CDClk (CDClk);
      Config.CDClk := CDClk;
      Config.Max_CDClk := CDClk;

      Get_Raw_Clock (Config.Raw_Clock);
   end Initialize;

   procedure Limit_Dotclocks
     (Configs        : in out Pipe_Configs;
      CDClk_Switch   :    out Boolean)
   is
   begin
      Config_Helpers.Limit_Dotclocks (Configs, Config.CDClk * 90 / 100);
      CDClk_Switch := False;
   end Limit_Dotclocks;

end HW.GFX.GMA.Power_And_Clocks;
