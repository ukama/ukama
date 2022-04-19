--
-- Copyright (C) 2014-2016, 2019 secunet Security Networks AG
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

with HW.Debug;
with HW.GFX.GMA.Registers;

use HW.GFX.GMA.Registers;

package body HW.GFX.GMA.PCode is

   GT_MAILBOX_READY : constant := 1 * 2 ** 31;

   -- Send a command and optionally wait for and return the reply.
   procedure Mailbox_Write_Read
     (MBox        : in     Word32;
      Command     : in     Word64;
      Reply       :    out Word64;
      Wait_Ready  : in     Boolean := False;
      Wait_Ack    : in     Boolean := True;
      Success     :    out Boolean)
   with
      Pre => Mailbox_Ready or Wait_Ready,
      Post => (if Wait_Ack and Success then Mailbox_Ready)
   is
      use type HW.Word64;

      Data : Word32;
   begin
      pragma Debug (Debug.Put_Line (GNAT.Source_Info.Enclosing_Entity));

      Reply := 0;
      Success := True;

      if Wait_Ready then
         Wait_Unset_Mask (GT_MAILBOX, GT_MAILBOX_READY, Success => Success);
         if not Success then
            return;
         end if;
      end if;

      Write (GT_MAILBOX_DATA, Word32 (Command and 16#ffff_ffff#));
      Write (GT_MAILBOX_DATA_1, Word32 (Shift_Right (Command, 32)));
      Write (GT_MAILBOX, GT_MAILBOX_READY or MBox);
      Mailbox_Ready := False;

      if Wait_Ack then
         Wait_Unset_Mask (GT_MAILBOX, GT_MAILBOX_READY, Success => Success);
         Mailbox_Ready := Success;

         Read (GT_MAILBOX_DATA, Data);
         Reply := Word64 (Data);
         Read (GT_MAILBOX_DATA_1, Data);
         Reply := Shift_Left (Word64 (Data), 32) or Reply;

         Write (GT_MAILBOX_DATA, 0);
         Write (GT_MAILBOX_DATA_1, 0);
      end if;
   end Mailbox_Write_Read;

   procedure Mailbox_Write
     (MBox        : in     Word32;
      Command     : in     Word64;
      Wait_Ready  : in     Boolean := False;
      Wait_Ack    : in     Boolean := True;
      Success     :    out Boolean)
   is
      pragma Warnings (GNATprove, Off, "unused assignment to ""Ignored_R""");
      Ignored_R : Word64;
   begin
      Mailbox_Write_Read
        (MBox, Command, Ignored_R, Wait_Ready, Wait_Ack, Success);
   end Mailbox_Write;

   procedure Mailbox_Request
     (MBox        : in     Word32;
      Command     : in     Word64;
      Reply_Mask  : in     Word64;
      Reply       : in     Word64 := 16#ffff_ffff_ffff_ffff#;
      TOut_MS     : in     Natural := Registers.Default_Timeout_MS;
      Wait_Ready  : in     Boolean := False;
      Success     :    out Boolean)
   is
      use type HW.Word64;

      Timeout : constant Time.T := Time.MS_From_Now (TOut_MS);
      Timed_Out : Boolean := False;

      Received_Reply : Word64;
   begin
      Success := False;
      loop
         pragma Loop_Invariant ((not Success and Wait_Ready) or Mailbox_Ready);
         Mailbox_Write_Read
           (MBox        => MBox,
            Command     => Command,
            Reply       => Received_Reply,
            Wait_Ready  => not Success and Wait_Ready,
            Success     => Success);
         exit when not Success;

         if (Received_Reply and Reply_Mask) = (Reply and Reply_Mask) then
            -- Ignore timeout if we succeeded anyway.
            Timed_Out := False;
            exit;
         end if;
         exit when Timed_Out;

         Timed_Out := Time.Timed_Out (Timeout);
      end loop;

      Success := Success and then not Timed_Out;
   end Mailbox_Request;

   procedure Mailbox_Write
     (MBox        : Word32;
      Command     : Word64;
      Wait_Ready  : Boolean := False)
   is
      pragma Warnings (GNATprove, Off, "unused assignment to ""Ignored_S""");
      Ignored_S : Boolean;
   begin
      Mailbox_Write (MBox, Command, Wait_Ready, False, Ignored_S);
   end Mailbox_Write;

end HW.GFX.GMA.PCode;
