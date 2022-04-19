--
-- Copyright (C) 2016, 2019 secunet Security Networks AG
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

with HW.GFX.GMA.Registers;

private package HW.GFX.GMA.PCode is

   -- We have to ensure that previous usage of the mailbox finished
   -- (Wait_Ready) or know that we already did so (Mailbox_Ready).
   --
   -- If we wait for the other side to acknowledge (Wait_Ack), we
   -- know that it's ready (=> Mailbox_Ready).

   -- XXX: Supposed to be a `Ghost` variable, but GNAT seems too broken?
   Mailbox_Ready : Boolean with Part_Of => HW.GFX.GMA.State;

   -- Just send a command, discard the reply.
   procedure Mailbox_Write
     (MBox        : in     Word32;
      Command     : in     Word64;
      Wait_Ready  : in     Boolean := False;
      Wait_Ack    : in     Boolean := True;
      Success     :    out Boolean)
   with
      Pre => Mailbox_Ready or Wait_Ready,
      Post => (if Wait_Ack and Success then Mailbox_Ready);

   -- Repeatedly send a request command the expected reply is received.
   procedure Mailbox_Request
     (MBox        : in     Word32;
      Command     : in     Word64;
      Reply_Mask  : in     Word64;
      Reply       : in     Word64 := 16#ffff_ffff_ffff_ffff#;
      TOut_MS     : in     Natural := Registers.Default_Timeout_MS;
      Wait_Ready  : in     Boolean := False;
      Success     :    out Boolean)
   with
      Pre => Mailbox_Ready or Wait_Ready,
      Post => (if Success then Mailbox_Ready);

   -- For final mailbox commands that don't have to wait.
   procedure Mailbox_Write
     (MBox        : Word32;
      Command     : Word64;
      Wait_Ready  : Boolean := False)
   with
      Pre => Mailbox_Ready or Wait_Ready;

end HW.GFX.GMA.PCode;
