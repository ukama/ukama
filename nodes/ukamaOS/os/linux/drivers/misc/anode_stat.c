 #include <linux/err.h>
 #include <linux/module.h>
 #include <linux/types.h>
 #include <linux/platform_device.h>
 #include <linux/io.h>
 #include <linux/of_gpio.h>
 #include <linux/of.h>
 #include <linux/sysfs.h>

struct ocfema_stat_priv
{
    	struct device *dev;
	
	int rf_therm_trip_gpio;
	char rf_therm_trip_szGpio[32];
	u8 rf_therm_trip_val;
	u32 rf_therm_trip_default;
	
	int rf_therm_alert_gpio;
	char rf_therm_alert_szGpio[32];
	u8 rf_therm_alert_val;
	u32 rf_therm_alert_default;

	int pg_reg_5p_7v_gpio;
	char pg_reg_5p_7v_szGpio[32];
	u8 pg_reg_5p_7v_val;
	u32 pg_reg_5p_7v_default;
			
		
	int pg_ldo_3p3_gpio;
	char pg_ldo_3p3_szGpio[32];
	u8 pg_ldo_3p3_val;
	u32 pg_ldo_3p3_default;

	int pg_ldo_5v_gpio;
	char pg_ldo_5v_szGpio[32];
	u8 pg_ldo_5v_val;
	u32 pg_ldo_5v_default;

	int rf_temp_alert_gpio;
	char rf_temp_alert_szGpio[32];
	u8 rf_temp_alert_val;
	u32 rf_temp_alert_default;
/*
	int config_rst_sw_gpio;
	char config_rst_sw_szGpio[32];
	u8 config_rst_sw_val;
	u32 config_rst_sw_default;
*/
};

static ssize_t show_rf_therm_trip(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_stat_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "RF thermal trip is in %s state.", ((priv->rf_therm_trip_val)==0)?"active":"inactive");
	return sprintf(buf, "%d\n", priv->rf_therm_trip_val);
}

static ssize_t show_rf_therm_alert(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_stat_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "RF thermal alert is in %s state.", ((priv->rf_therm_alert_val)==0)?"active":"inactive");
	return sprintf(buf, "%d\n", priv->rf_therm_alert_val);
}

static ssize_t show_pg_reg_5p_7V(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_stat_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "Power Good status for 5.7v is %s.", ((priv->pg_reg_5p_7v_val)==0)?"ok":"not ok");
	return sprintf(buf, "%d\n", priv->pg_reg_5p_7v_val);
}

static ssize_t show_pg_ldo_3p3(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_stat_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "Power Good status for 3.3V LDO is %s.", ((priv->pg_ldo_3p3_val)==0)?"ok":"not ok");
	return sprintf(buf, "%d\n", priv->pg_ldo_3p3_val);
}

static ssize_t show_pg_ldo_5V(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_stat_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "Power Good status for 5V LDO is %s.", ((priv->pg_ldo_5v_val)==0)?"ok":"not ok");
	return sprintf(buf, "%d\n", priv->pg_ldo_5v_val);
}

static ssize_t show_rf_temp_alert(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_stat_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "RF temperature alert is in %s state.", ((priv->rf_temp_alert_val)==0)?"active":"inactive");
	return sprintf(buf, "%d\n", priv->rf_temp_alert_val);
}
/*
static ssize_t show_config_rst_sw(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_stat_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "Reset switch is %s.", ((priv->config_rst_sw_val)==0)?"asserted":"not asserted");
	return sprintf(buf, "%d\n", priv->config_rst_sw_val);
}
*/
static DEVICE_ATTR(rf_therm_trip, S_IRUGO, show_rf_therm_trip, NULL);
static DEVICE_ATTR(rf_therm_alert, S_IRUGO, show_rf_therm_alert, NULL);
static DEVICE_ATTR(pg_reg_5_7V, S_IRUGO,show_pg_reg_5p_7V, NULL);
static DEVICE_ATTR(pg_ldo_3p3V, S_IRUGO, show_pg_ldo_3p3, NULL);
static DEVICE_ATTR(pg_ldo_5V, S_IRUGO, show_pg_ldo_5V, NULL);
static DEVICE_ATTR(rf_temp_alert, S_IRUGO, show_rf_temp_alert, NULL);
//static DEVICE_ATTR(config_rst_sw, S_IRUGO, show_config_rst_sw, NULL);

static struct attribute *ocfema_stat_attrs[] = {
	&dev_attr_rf_therm_trip.attr,
	&dev_attr_rf_therm_alert.attr,
	&dev_attr_pg_reg_5_7V.attr,
	&dev_attr_pg_ldo_3p3V.attr,
	&dev_attr_pg_ldo_5V.attr,
	&dev_attr_rf_temp_alert.attr,
	//&dev_attr_config_rst_sw.attr,
	NULL,
};

static const struct attribute_group ocfema_stat_attr_group = {
	.attrs = ocfema_stat_attrs,
};

int ocfema_stat_parse_dt(struct platform_device *pdev)
{
	int ret = 0; 
	struct device_node *np = pdev->dev.of_node;
	struct ocfema_stat_priv *priv = dev_get_drvdata(&pdev->dev);
	if (!np) {
		return -EINVAL;
	}
	
	priv->rf_therm_trip_gpio = of_get_named_gpio(np, "rf-therm-trip-gpio", 0);
        if (priv->rf_therm_trip_gpio < 0) {
                dev_warn(&pdev->dev, "Can't read RF thermal trip gpio from device tree.\n");
        }

	priv->rf_therm_alert_gpio = of_get_named_gpio(np, "therm-alert-gpio", 0);
        if (priv->rf_therm_alert_gpio < 0) {
                dev_warn(&pdev->dev, "Can't read RF thermal alert gpio from device tree.\n");
        }
	
	priv->pg_reg_5p_7v_gpio = of_get_named_gpio(np, "pg-reg-5p-7v-gpio", 0);
        if (priv->pg_reg_5p_7v_gpio < 0) {
                dev_warn(&pdev->dev, "Can't read power good gpio for 5.7V from device tree.\n");
        }
	
	priv->pg_ldo_3p3_gpio = of_get_named_gpio(np, "pg-ldo-3p3-gpio", 0);
        if (priv->pg_ldo_3p3_gpio < 0) {
                dev_warn(&pdev->dev, "Can't read power good for 3.3V LDO from device tree.\n");
        }
	
	priv->pg_ldo_5v_gpio = of_get_named_gpio(np, "pg-ldo-5v-gpio", 0);
        if (priv->pg_ldo_5v_gpio < 0) {
                dev_warn(&pdev->dev, "Can't read power good for 5V LDO from device tree.\n");
        }

	priv->rf_temp_alert_gpio = of_get_named_gpio(np, "rf-temp-alert-gpio", 0);
        if (priv->rf_temp_alert_gpio < 0) {
                dev_warn(&pdev->dev, "Can't read RF temp alert gpio from device tree.\n");
	}
/*	
	priv->config_rst_sw_gpio = of_get_named_gpio(np, "config-rst-sw-gpio", 0);
        if (priv->config_rst_sw_gpio < 0) {
                dev_warn(&pdev->dev, "Can't read reset switch gpio from device tree.\n");
	}
*/	
	return ret;
}

static int ocfema_stat_probe(struct platform_device *pdev)
{ 
	struct ocfema_stat_priv *priv;
	int ret;
	
    	priv = devm_kzalloc(&pdev->dev, sizeof(struct ocfema_stat_priv), GFP_KERNEL);
	if (!priv) {
		dev_err(priv->dev, "Unable to allocate memory.\n");
		return -ENOMEM;
	}

	priv->dev = &pdev->dev;
	platform_set_drvdata(pdev, priv);

	ret = ocfema_stat_parse_dt(pdev);
	if (ret) {	
		dev_err(&pdev->dev, "Parsing failed for ocfema status node.");
		return ret;
	}	
	
	dev_info(&pdev->dev, "Configure OC-FEMA status pins on boot.");
   	
	sprintf(priv->rf_therm_trip_szGpio, "rf_therm_trip");
	ret = gpio_request(priv->rf_therm_trip_gpio, priv->rf_therm_trip_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->rf_therm_trip_gpio, priv->rf_therm_trip_szGpio);
		return ret;
	}

	ret = gpio_direction_input(priv->rf_therm_trip_gpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->rf_therm_trip_gpio, ret);
		return ret;
	}
	
	dev_info(&pdev->dev, "RF therm trip configured.");
   
	sprintf(priv->rf_therm_alert_szGpio, "rf_therm_alert");
	ret = gpio_request(priv->rf_therm_alert_gpio, priv->rf_therm_alert_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->rf_therm_alert_gpio, priv->rf_therm_alert_szGpio);
		return ret;
	}

	ret = gpio_direction_input(priv->rf_therm_alert_gpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->rf_therm_alert_gpio, ret);
		return ret;
	}
	
	dev_info(&pdev->dev, "RF therm alert configured.");
   	
	sprintf(priv->pg_reg_5p_7v_szGpio, "pg_reg_5p_7V");
	ret = gpio_request(priv->pg_reg_5p_7v_gpio, priv->pg_reg_5p_7v_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->pg_reg_5p_7v_gpio, priv->pg_reg_5p_7v_szGpio);
		return ret;
	}

	ret = gpio_direction_input(priv->pg_reg_5p_7v_gpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->pg_reg_5p_7v_gpio, ret);
		return ret;
	}
	
	dev_info(&pdev->dev, "Power good for 5.7V is configured.");
   	
	sprintf(priv->pg_ldo_3p3_szGpio, "pg_ldo_3p3");
	ret = gpio_request(priv->pg_ldo_3p3_gpio, priv->pg_ldo_3p3_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->pg_ldo_3p3_gpio, priv->pg_ldo_3p3_szGpio);
		return ret;
	}

	ret = gpio_direction_input(priv->pg_ldo_3p3_gpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->pg_ldo_3p3_gpio, ret);
		return ret;
	}
	
	dev_info(&pdev->dev, "Power good for 3.3V LDO is configured.");
	
			
	sprintf(priv->pg_ldo_5v_szGpio, "pg_ldo_5V");
	ret = gpio_request(priv->pg_ldo_5v_gpio, priv->pg_ldo_5v_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->pg_ldo_5v_gpio, priv->pg_ldo_5v_szGpio);
		return ret;
	}

	ret = gpio_direction_input(priv->pg_ldo_5v_gpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->pg_ldo_5v_gpio, ret);
		return ret;
	}
	
	dev_info(&pdev->dev, "Power good for 5V LDO is configured.");
	
		
	sprintf(priv->rf_temp_alert_szGpio, "rf_temp_alert");
	ret = gpio_request(priv->rf_temp_alert_gpio, priv->rf_temp_alert_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->rf_temp_alert_gpio, priv->rf_temp_alert_szGpio);
		return ret;
	}

	ret = gpio_direction_input(priv->rf_temp_alert_gpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->rf_temp_alert_gpio, ret);
		return ret;
	}
	
	dev_info(&pdev->dev, "RF temp alert configured.");
/*	
	sprintf(priv->config_rst_sw_szGpio, "config_rst_sw");
	ret = gpio_request(priv->config_rst_sw_gpio, priv->config_rst_sw_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->config_rst_sw_gpio, priv->config_rst_sw_szGpio);
		return ret;
	}

	ret = gpio_direction_input(priv->config_rst_sw_gpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->config_rst_sw_gpio, ret);
		return ret;
	}
	
	dev_info(&pdev->dev, "Reset switch is configured.");
*/	
	dev_info(&pdev->dev, "Setting Sysfs for OC-FEMA Status.");
	
	ret = sysfs_create_group(&priv->dev->kobj, &ocfema_stat_attr_group);
	if (ret) {
		dev_err(priv->dev, "unable to create sysfs files\n");
		return ret;
	}
	return 0;
}

static int ocfema_stat_remove(struct platform_device *pdev)
{
	struct ocfema_stat_priv *priv = platform_get_drvdata(pdev);
	sysfs_remove_group(&priv->dev->kobj, &ocfema_stat_attr_group);
	return 0;
}

static const struct of_device_id ocfema_of_match[] = {
	{ .compatible = "oc,fema-stat", },
	{}
};
MODULE_DEVICE_TABLE(of, ocfema_of_match);

static struct platform_driver ocfema_stat_driver = {
	.driver = {
		.name = "ocfema_stat",
		.owner = THIS_MODULE,
		.of_match_table = of_match_ptr(ocfema_of_match),
	},
	.probe = ocfema_stat_probe,
	.remove = ocfema_stat_remove,
};

module_platform_driver(ocfema_stat_driver);

MODULE_DESCRIPTION("OCFEMA Status Driver");
MODULE_AUTHOR("<vthakur@fb.com>");
MODULE_LICENSE("GPL");
