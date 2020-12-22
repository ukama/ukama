#include "storage_interface.h"

int Read_Module_Count(uint8_t *module_count)
{
    /* todo: Replace with call to the function that will be provided by Vishal. */
    *module_count = 2;

    return STORAGE_INTERFACE_SUCCESS;
}

int Read_Each_Module_UUID(char **uuid, uint8_t module_number)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;
    uint8_t module_count = -1;

    /* If given 'module_number' is valid and UUID is not null. */
    if ((uuid) && (Read_Module_Count(&module_count) == STORAGE_INTERFACE_SUCCESS) && (module_number < module_count))
    {
        /* todo: Replace with BSP function that will be provided by Vishal. */
        char *temp = "213452547";

        /* todo: Might have to deep copy this. */
        *uuid = temp;

        status = STORAGE_INTERFACE_SUCCESS;
    }

    return (status);
}

int Read_Each_Module_Manufacturer(char **manufacturer, uint8_t module_number)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;
    uint8_t module_count = -1;

    /* If given 'module_number' is valid and manufacturer is not null */
    if ((manufacturer) && (Read_Module_Count(&module_count) == STORAGE_INTERFACE_SUCCESS) && (module_number < module_count))
    {
        /* todo: Replace with BSP function that will be provided by Vishal. */
        char *temp = "Ukama Inc";

        /* todo: Might have to deep copy this. */
        *manufacturer = temp;

        status = STORAGE_INTERFACE_SUCCESS;
    }

    return (status);
}

int Read_Each_Module_Model(char **model, uint8_t module_number)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;
    uint8_t module_count = -1;

    /* If given 'model' is valid and manufacturer is not null */
    if ((model) && (Read_Module_Count(&module_count) == STORAGE_INTERFACE_SUCCESS) && (module_number < module_count))
    {
        /* todo: Replace with BSP function that will be provided by Vishal. */
        char *temp = "COM_x86_4x1.9GHz_2GB_V2.0_Rev-A";

        /* todo: Might have to deep copy this. */
        *model = temp;

        status = STORAGE_INTERFACE_SUCCESS;
    }

    return (status);
}

int Read_Each_Module_PartNumber(char **partnum, uint8_t module_number)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;
    uint8_t module_count = -1;

    /* If given 'partnum' is valid and manufacturer is not null */
    if ((partnum) && (Read_Module_Count(&module_count) == STORAGE_INTERFACE_SUCCESS) && (module_number < module_count))
    {
        /* todo: Replace with BSP function that will be provided by Vishal. */
        char *temp = "AC102J";

        /* todo: Might have to deep copy this. */
        *partnum = temp;

        status = STORAGE_INTERFACE_SUCCESS;
    }

    return (status);
}

int Read_Each_Module_Mfgdate(char **mfgdate, uint8_t module_number)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;
    uint8_t module_count = -1;

    /* If given 'mfgdate' is valid and manufacturer is not null */
    if ((mfgdate) && (Read_Module_Count(&module_count) == STORAGE_INTERFACE_SUCCESS) && (module_number < module_count))
    {
        /* todo: Replace with BSP function that will be provided by Vishal. */
        char *temp = "16-july-2018";

        /* todo: Might have to deep copy this. */
        *mfgdate = temp;

        status = STORAGE_INTERFACE_SUCCESS;
    }

    return (status);
}

int Read_Each_Module_Moduleclass(char **moduleclass, uint8_t module_number)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;
    uint8_t module_count = -1;

    /* If given 'moduleclass' is valid and manufacturer is not null */
    if ((moduleclass) && (Read_Module_Count(&module_count) == STORAGE_INTERFACE_SUCCESS) && (module_number < module_count))
    {
        /* todo: Replace with BSP function that will be provided by Vishal. */
        char *temp = "Radio";

        /* todo: Might have to deep copy this. */
        *moduleclass = temp;

        status = STORAGE_INTERFACE_SUCCESS;
    }

    return (status);
}

int Read_Each_Module_Swversion(char **swversion, uint8_t module_number)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;
    uint8_t module_count = -1;

    /* If given 'swversion' is valid and manufacturer is not null */
    if ((swversion) && (Read_Module_Count(&module_count) == STORAGE_INTERFACE_SUCCESS) && (module_number < module_count))
    {
        /* todo: Replace with BSP function that will be provided by Vishal. */
        char *temp = "10.2";

        /* todo: Might have to deep copy this. */
        *swversion = temp;

        status = STORAGE_INTERFACE_SUCCESS;
    }

    return (status);
}

int Read_Each_Module_Hwversion(char **hwversion, uint8_t module_number)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;
    uint8_t module_count = -1;

    /* If given 'hwversion' is valid and manufacturer is not null */
    if ((hwversion) && (Read_Module_Count(&module_count) == STORAGE_INTERFACE_SUCCESS) && (module_number < module_count))
    {
        /* todo: Replace with BSP function that will be provided by Vishal. */
        char *temp = "1.0";

        /* todo: Might have to deep copy this. */
        *hwversion = temp;

        status = STORAGE_INTERFACE_SUCCESS;
    }

    return (status);
}

int Read_Temperature_Sensor_Count(char *model_class, uint8_t *count)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;

    if (model_class != NULL && count != NULL)
    {
        status = STORAGE_INTERFACE_SUCCESS;

        if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_MASK) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_TEMP_SENSORS_MASK;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_RADIO) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_TEMP_SENSORS_RADIO;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_COM) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_TEMP_SENSORS_COM;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_POWER) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_TEMP_SENSORS_POWER;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_CTRL) == 0)
        {
            /* As per the device management doc it only has one temperature sensore. */
            *count = STORAGE_INTERFACE_TEMP_SENSORS_CTRL;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_FEM) == 0)
        {
            /* As per the device management doc it only has one temperature sensore. */
            *count = STORAGE_INTERFACE_TEMP_SENSORS_FEM;
        }
        else
        {
            status = STORAGE_INTERFACE_UNKNOWN_MODULE_CLASS;
        }
    }

    return status;
}

int Read_Digital_Input_Count(char *model_class, uint8_t *count)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;

    if (model_class != NULL && count != NULL)
    {
        status = STORAGE_INTERFACE_SUCCESS;

        if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_MASK) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_DIGITAL_INPUT_MASK;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_RADIO) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_DIGITAL_INPUT_RADIO;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_COM) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_DIGITAL_INPUT_COM;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_POWER) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_DIGITAL_INPUT_POWER;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_CTRL) == 0)
        {
            /* As per the device management doc it only has one temperature sensore. */
            *count = STORAGE_INTERFACE_DIGITAL_INPUT_CTRL;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_FEM) == 0)
        {
            /* As per the device management doc it only has one temperature sensore. */
            *count = STORAGE_INTERFACE_DIGITAL_INPUT_FEM;
        }
        else
        {
            status = STORAGE_INTERFACE_UNKNOWN_MODULE_CLASS;
        }
    }

    return status;
}

int Read_Digital_Output_Count(char *model_class, uint8_t *count)
{
    int status = STORAGE_INTERFACE_INVALID_PARAM;

    if (model_class != NULL && count != NULL)
    {
        status = STORAGE_INTERFACE_SUCCESS;

        if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_MASK) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_DIGITAL_OUTPUT_MASK;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_RADIO) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_DIGITAL_OUTPUT_RADIO;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_COM) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_DIGITAL_OUTPUT_COM;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_POWER) == 0)
        {
            //todo: the device management doc does not have any information
            *count = STORAGE_INTERFACE_DIGITAL_OUTPUT_POWER;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_CTRL) == 0)
        {
            /* As per the device management doc it only has one temperature sensore. */
            *count = STORAGE_INTERFACE_DIGITAL_OUTPUT_CTRL;
        }
        else if (strcmp(model_class, STORAGE_INTERFACE_MODULE_CLASS_FEM) == 0)
        {
            /* As per the device management doc it only has one temperature sensore. */
            *count = STORAGE_INTERFACE_DIGITAL_OUTPUT_FEM;
        }
        else
        {
            status = STORAGE_INTERFACE_UNKNOWN_MODULE_CLASS;
        }
    }

    return status;
}