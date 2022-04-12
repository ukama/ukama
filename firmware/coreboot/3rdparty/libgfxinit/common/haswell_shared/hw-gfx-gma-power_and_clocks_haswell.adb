--
-- Copyright (C) 2014-2018 secunet Security Networks AG
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

with GNAT.Source_Info;

with HW.Time;
with HW.Debug;
with HW.GFX.GMA.Config;
with HW.GFX.GMA.PCode;
with HW.GFX.GMA.Registers;

package body HW.GFX.GMA.Power_And_Clocks_Haswell is

   LCPLL_CTL_CD_FREQ_SEL_MASK          : constant := 3 * 2 ** 26;
   LCPLL_CTL_CD_FREQ_SEL_450_MHZ       : constant := 0 * 2 ** 26;
   LCPLL_CTL_CD_FREQ_SEL_HSW_ALTERNATE : constant := 1 * 2 ** 26;
   LCPLL_CTL_CD_FREQ_SEL_BDW_540_MHZ   : constant := 1 * 2 ** 26;
   LCPLL_CTL_CD_FREQ_SEL_BDW_337_5_MHZ : constant := 2 * 2 ** 26;
   LCPLL_CTL_CD_FREQ_SEL_BDW_675_MHZ   : constant := 3 * 2 ** 26;
   LCPLL_CTL_CD_SOURCE_SELECT_FCLK     : constant := 1 * 2 ** 21;
   LCPLL_CTL_CD_SOURCE_FCLK_DONE       : constant := 1 * 2 ** 19;

   function LCPLL_CTL_CD_FREQ_SEL_BDW (CDClk : Config.CDClk_Range) return Word32
   is
     (case CDClk is
         when 675_000_000 => LCPLL_CTL_CD_FREQ_SEL_BDW_675_MHZ,
         when 540_000_000 => LCPLL_CTL_CD_FREQ_SEL_BDW_540_MHZ,
         when 450_000_000 => LCPLL_CTL_CD_FREQ_SEL_450_MHZ,
         when others      => LCPLL_CTL_CD_FREQ_SEL_BDW_337_5_MHZ);

   FUSE_STRAP_DISPLAY_CDCLK_LIMIT      : constant := 1 * 2 ** 24;

   HSW_PCODE_DE_WRITE_FREQ             : constant := 16#17#;
   BDW_PCODE_DISPLAY_FREQ_CHANGE       : constant := 16#18#;

   ----------------------------------------------------------------------------

   PWR_WELL_CTL_ENABLE_REQUEST   : constant := 1 * 2 ** 31;
   PWR_WELL_CTL_DISABLE_REQUEST  : constant := 0 * 2 ** 31;
   PWR_WELL_CTL_STATE_ENABLED    : constant := 1 * 2 ** 30;

   ----------------------------------------------------------------------------

   SRD_CTL_ENABLE          : constant := 1 * 2 ** 31;
   SRD_STATUS_STATE_MASK   : constant := 7 * 2 ** 29;

   type Pipe is (EDP, A, B, C);
   type SRD_Regs is record
      CTL     : Registers.Registers_Index;
      STATUS  : Registers.Registers_Index;
   end record;
   type SRD_Per_Pipe_Regs is array (Pipe) of SRD_Regs;
   SRD : constant SRD_Per_Pipe_Regs := SRD_Per_Pipe_Regs'
     (A     => SRD_Regs'
        (CTL      => Registers.SRD_CTL_A,
         STATUS   => Registers.SRD_STATUS_A),
      B     => SRD_Regs'
        (CTL      => Registers.SRD_CTL_B,
         STATUS   => Registers.SRD_STATUS_B),
      C     => SRD_Regs'
        (CTL      => Registers.SRD_CTL_C,
         STATUS   => Registers.SRD_STATUS_C),
      EDP   => SRD_Regs'
        (CTL      => Registers.SRD_CTL_EDP,
         STATUS   => Registers.SRD_STATUS_EDP));

   ----------------------------------------------------------------------------

   IPS_CTL_ENABLE          : constant := 1 * 2 ** 31;
   DISPLAY_IPS_CONTROL     : constant := 16#19#;

   ----------------------------------------------------------------------------

   procedure PSR_Off
   is
      Enabled : Boolean;
   begin
      pragma Debug (Debug.Put_Line (GNAT.Source_Info.Enclosing_Entity));

      if Config.Has_Per_Pipe_SRD then
         for P in Pipe loop
            Registers.Is_Set_Mask (SRD (P).CTL, SRD_CTL_ENABLE, Enabled);
            if Enabled then
               Registers.Unset_Mask (SRD (P).CTL, SRD_CTL_ENABLE);
               Registers.Wait_Unset_Mask (SRD (P).STATUS, SRD_STATUS_STATE_MASK);

               pragma Debug (Debug.Put_Line ("Disabled PSR."));
            end if;
         end loop;
      else
         Registers.Is_Set_Mask (Registers.SRD_CTL, SRD_CTL_ENABLE, Enabled);
         if Enabled then
            Registers.Unset_Mask (Registers.SRD_CTL, SRD_CTL_ENABLE);
            Registers.Wait_Unset_Mask (Registers.SRD_STATUS, SRD_STATUS_STATE_MASK);

            pragma Debug (Debug.Put_Line ("Disabled PSR."));
         end if;
      end if;
   end PSR_Off;

   ----------------------------------------------------------------------------

   procedure IPS_Off
   is
      Enabled : Boolean;
   begin
      pragma Debug (Debug.Put_Line (GNAT.Source_Info.Enclosing_Entity));

      if Config.Has_IPS then
         Registers.Is_Set_Mask (Registers.IPS_CTL, IPS_CTL_ENABLE, Enabled);
         if Enabled then
            if Config.Has_IPS_CTL_Mailbox then
               PCode.Mailbox_Write (DISPLAY_IPS_CONTROL, 0, Wait_Ready => True);
               Registers.Wait_Unset_Mask
                 (Register => Registers.IPS_CTL,
                  Mask     => IPS_CTL_ENABLE,
                  TOut_MS  => 42);
            else
               Registers.Unset_Mask (Registers.IPS_CTL, IPS_CTL_ENABLE);
            end if;

            pragma Debug (Debug.Put_Line ("Disabled IPS."));
            -- We have to wait until the next vblank here.
            -- 20ms should be enough.
            Time.M_Delay (20);
         end if;
      end if;
   end IPS_Off;

   ----------------------------------------------------------------------------

   procedure PDW_Off
   is
      Ctl1, Ctl2, Ctl3, Ctl4 : Word32;
   begin
      pragma Debug (Debug.Put_Line (GNAT.Source_Info.Enclosing_Entity));

      Registers.Read (Registers.PWR_WELL_CTL_BIOS, Ctl1);
      Registers.Read (Registers.PWR_WELL_CTL_DRIVER, Ctl2);
      Registers.Read (Registers.PWR_WELL_CTL_KVMR, Ctl3);
      Registers.Read (Registers.PWR_WELL_CTL_DEBUG, Ctl4);
      pragma Debug (Registers.Posting_Read (Registers.PWR_WELL_CTL5)); --  Result for debugging only
      pragma Debug (Registers.Posting_Read (Registers.PWR_WELL_CTL6)); --  Result for debugging only

      if ((Ctl1 or Ctl2 or Ctl3 or Ctl4) and
          PWR_WELL_CTL_ENABLE_REQUEST) /= 0
      then
         Registers.Wait_Set_Mask
           (Registers.PWR_WELL_CTL_DRIVER, PWR_WELL_CTL_STATE_ENABLED);
      end if;

      if (Ctl1 and PWR_WELL_CTL_ENABLE_REQUEST) /= 0 then
         Registers.Write (Registers.PWR_WELL_CTL_BIOS, PWR_WELL_CTL_DISABLE_REQUEST);
      end if;

      if (Ctl2 and PWR_WELL_CTL_ENABLE_REQUEST) /= 0 then
         Registers.Write (Registers.PWR_WELL_CTL_DRIVER, PWR_WELL_CTL_DISABLE_REQUEST);
      end if;
   end PDW_Off;

   procedure PDW_On
   is
      Ctl1, Ctl2, Ctl3, Ctl4 : Word32;
   begin
      pragma Debug (Debug.Put_Line (GNAT.Source_Info.Enclosing_Entity));

      Registers.Read (Registers.PWR_WELL_CTL_BIOS, Ctl1);
      Registers.Read (Registers.PWR_WELL_CTL_DRIVER, Ctl2);
      Registers.Read (Registers.PWR_WELL_CTL_KVMR, Ctl3);
      Registers.Read (Registers.PWR_WELL_CTL_DEBUG, Ctl4);
      pragma Debug (Registers.Posting_Read (Registers.PWR_WELL_CTL5)); --  Result for debugging only
      pragma Debug (Registers.Posting_Read (Registers.PWR_WELL_CTL6)); --  Result for debugging only

      if ((Ctl1 or Ctl2 or Ctl3 or Ctl4) and
          PWR_WELL_CTL_ENABLE_REQUEST) = 0
      then
         Registers.Wait_Unset_Mask
           (Registers.PWR_WELL_CTL_DRIVER, PWR_WELL_CTL_STATE_ENABLED);
      end if;

      if (Ctl2 and PWR_WELL_CTL_ENABLE_REQUEST) = 0 then
         Registers.Write (Registers.PWR_WELL_CTL_DRIVER, PWR_WELL_CTL_ENABLE_REQUEST);
         Registers.Wait_Set_Mask
           (Registers.PWR_WELL_CTL_DRIVER, PWR_WELL_CTL_STATE_ENABLED);
      end if;
   end PDW_On;

   function Need_PDW (Checked_Configs : Pipe_Configs) return Boolean
   is
      Primary : Pipe_Config renames Checked_Configs (GMA.Primary);
   begin
      return
         (Config.Use_PDW_For_EDP_Scaling and then
          (Primary.Port = Internal and Requires_Scaling (Primary)))
         or
         (Primary.Port /= Disabled and Primary.Port /= Internal)
         or
         Checked_Configs (Secondary).Port /= Disabled
         or
         Checked_Configs (Tertiary).Port /= Disabled;
   end Need_PDW;

   ----------------------------------------------------------------------------

   procedure Pre_All_Off is
   begin
      -- HSW: disable panel self refresh (PSR) on eDP if enabled
         -- wait for PSR idling
      PSR_Off;
      IPS_Off;
   end Pre_All_Off;

   function Normalize_CDClk (CDClk : in Int64) return Config.CDClk_Range is
     (   if CDClk <= 337_500_000 then 337_500_000
      elsif CDClk <= 450_000_000 then 450_000_000
      elsif CDClk <= 540_000_000 then 540_000_000
                                 else 675_000_000);

   procedure Get_Cur_CDClk (CDClk : out Config.CDClk_Range)
   is
      LCPLL_CTL : Word32;
   begin
      Registers.Read (Registers.LCPLL_CTL, LCPLL_CTL);
      CDClk :=
        (if Config.Has_Broadwell_CDClk then
           (case LCPLL_CTL and LCPLL_CTL_CD_FREQ_SEL_MASK is
               when LCPLL_CTL_CD_FREQ_SEL_BDW_540_MHZ    => 540_000_000,
               when LCPLL_CTL_CD_FREQ_SEL_BDW_337_5_MHZ  => 337_500_000,
               when LCPLL_CTL_CD_FREQ_SEL_BDW_675_MHZ    => 675_000_000,
               when others                               => 450_000_000)
         else
           (case LCPLL_CTL and LCPLL_CTL_CD_FREQ_SEL_MASK is
               when LCPLL_CTL_CD_FREQ_SEL_HSW_ALTERNATE =>
                 (if    Config.Is_ULX  then 337_500_000
                  elsif Config.Is_ULT  then 450_000_000
                                       else 540_000_000),
               when others => 450_000_000));
   end Get_Cur_CDClk;

   procedure Get_Max_CDClk (CDClk : out Config.CDClk_Range)
   is
      FUSE_STRAP : Word32;
   begin
      if Config.Has_Broadwell_CDClk then
         Registers.Read (Registers.FUSE_STRAP, FUSE_STRAP);
         CDClk :=
           (if (FUSE_STRAP and FUSE_STRAP_DISPLAY_CDCLK_LIMIT) /= 0 then
               450_000_000
            elsif Config.Is_ULX then
               450_000_000
            elsif Config.Is_ULT then
               540_000_000
            else
               675_000_000);
      else
         -- We may never switch CDClk on Haswell. So from our point
         -- of view, the CDClk we start with is the maximum.
         Get_Cur_CDClk (CDClk);
      end if;
   end Get_Max_CDClk;

   procedure Set_CDClk (CDClk_In : Frequency_Type)
   is
      CDClk : constant Config.CDClk_Range :=
         Normalize_CDClk (Frequency_Type'Min (CDClk_In, Config.Max_CDClk));
      Success : Boolean;
   begin
      if not Config.Can_Switch_CDClk then
         return;
      end if;

      PCode.Mailbox_Write
        (MBox        => BDW_PCODE_DISPLAY_FREQ_CHANGE,
         Command     => 0,
         Wait_Ready  => True,
         Success     => Success);

      if not Success then
         pragma Debug (Debug.Put_Line
           ("ERROR: PCODE didn't acknowledge frequency change."));
         return;
      end if;

      Registers.Set_Mask
        (Register => Registers.LCPLL_CTL,
         Mask     => LCPLL_CTL_CD_SOURCE_SELECT_FCLK);
      Registers.Wait_Set_Mask
        (Register => Registers.LCPLL_CTL,
         Mask     => LCPLL_CTL_CD_SOURCE_FCLK_DONE);

      Registers.Unset_And_Set_Mask
        (Register    => Registers.LCPLL_CTL,
         Mask_Unset  => LCPLL_CTL_CD_FREQ_SEL_MASK,
         Mask_Set    => LCPLL_CTL_CD_FREQ_SEL_BDW (CDClk));
      Registers.Posting_Read (Registers.LCPLL_CTL);

      Registers.Unset_Mask
        (Register => Registers.LCPLL_CTL,
         Mask     => LCPLL_CTL_CD_SOURCE_SELECT_FCLK);
      Registers.Wait_Unset_Mask
        (Register => Registers.LCPLL_CTL,
         Mask     => LCPLL_CTL_CD_SOURCE_FCLK_DONE);

      PCode.Mailbox_Write
        (MBox        => HSW_PCODE_DE_WRITE_FREQ,
         Command     => (case CDClk is
                           when 675_000_000 => 3,
                           when 540_000_000 => 1,
                           when 450_000_000 => 0,
                           when others      => 2));

      Registers.Write
        (Register => Registers.CDCLK_FREQ,
         Value    => Word32 (Div_Round_Closest (CDClk, 1_000_000) - 1));

      Config.CDClk := CDClk;
   end Set_CDClk;

   procedure Initialize
   is
      CDClk : Config.CDClk_Range;
   begin
      -- HSW: disable power down well
      PDW_Off;

      Get_Cur_CDClk (CDClk);
      Config.CDClk := CDClk;
      Get_Max_CDClk (CDClk);
      Config.Max_CDClk := CDClk;
      Set_CDClk (Config.Default_CDClk_Freq);

      Config.Raw_Clock := Config.Default_RawClk_Freq;
   end Initialize;

   procedure Limit_Dotclocks
     (Configs        : in out Pipe_Configs;
      CDClk_Switch   :    out Boolean)
   is
   begin
      Config_Helpers.Limit_Dotclocks (Configs, Config.Max_CDClk);
      CDClk_Switch :=
         Config.Can_Switch_CDClk and then
         Config.CDClk /= Normalize_CDClk
           (Config_Helpers.Highest_Dotclock (Configs));
   end Limit_Dotclocks;

   procedure Update_CDClk (Configs : in out Pipe_Configs)
   is
      New_CDClk : constant Frequency_Type :=
         Config_Helpers.Highest_Dotclock (Configs);
   begin
      Set_CDClk (New_CDClk);
      Config_Helpers.Limit_Dotclocks (Configs, Config.CDClk);
   end Update_CDClk;

   procedure Power_Set_To (Configs : Pipe_Configs) is
   begin
      if Need_PDW (Configs) then
         PDW_On;
      else
         PDW_Off;
      end if;
   end Power_Set_To;

   procedure Power_Up (Old_Configs, New_Configs : Pipe_Configs) is
   begin
      if not Need_PDW (Old_Configs) and Need_PDW (New_Configs) then
         PDW_On;
      end if;
   end Power_Up;

   procedure Power_Down (Old_Configs, Tmp_Configs, New_Configs : Pipe_Configs)
   is
   begin
      if (Need_PDW (Old_Configs) or Need_PDW (Tmp_Configs)) and
         not Need_PDW (New_Configs)
      then
         PDW_Off;
      end if;
   end Power_Down;

end HW.GFX.GMA.Power_And_Clocks_Haswell;
