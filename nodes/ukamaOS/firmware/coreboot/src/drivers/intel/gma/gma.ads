with Interfaces.C;

with HW.GFX.EDID;

package GMA is

   function read_edid
     (raw_edid :    out HW.GFX.EDID.Raw_EDID_Data;
      Port     : in     Interfaces.C.int)
      return Interfaces.C.int
   with
      Export, Convention => C, External_Name => "gma_read_edid";

end GMA;
