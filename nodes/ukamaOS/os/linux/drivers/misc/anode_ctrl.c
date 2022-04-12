 #include <linux/err.h>
 #include <linux/module.h>
 #include <linux/types.h>
 #include <linux/platform_device.h>
 #include <linux/io.h>
 #include <linux/of_gpio.h>
 #include <linux/of.h>
 #include <linux/sysfs.h>

struct ocfema_ctrl_priv
{
    	struct device *dev;
#if CHECK_WD	
	int som_wdi_in_gpio;
	char som_wdi_in_szGpio[32];
	u8 som_wdi_in_val;
	u32 som_wdi_in_default;
#endif	
	int rf_pwr_dis_gpio;
	char rf_pwr_dis_szGpio[32];
	u8 rf_pwr_dis_val;
	u32 rf_pwr_dis_default;
	
	int rf_eeprom_wp_en_gpio;
	char rf_eeprom_wp_en_szGpio[32];
	u8 rf_eeprom_wp_en_val;
	u32 rf_eeprom_wp_en_default;
	
	int pa_dis_gpio;
        char pa_dis_szGpio[32];
        u8 pa_dis_val;
        u32 pa_dis_default;

	int pga_pwr_dis_gpio;
	char pga_pwr_dis_szGpio[32];
	u8 pga_pwr_dis_val;
	u32 pga_pwr_dis_default;
			
		
	int mp_eeprom_wp_en_gpio;
	char mp_eeprom_wp_en_szGpio[32];
	u8 mp_eeprom_wp_en_val;
	u32 mp_eeprom_wp_en_default;
};

static ssize_t show_rf_pwr(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_ctrl_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "RF power is %s.", ((priv->rf_pwr_dis_val)==1)?"disabled":"enabled");
	return sprintf(buf, "%d\n", priv->rf_pwr_dis_val);
}

static ssize_t set_rf_pwr(struct device *dev,
		struct device_attribute *attr,
		const char *buf, size_t count)
{
	struct ocfema_ctrl_priv *priv = dev_get_drvdata(dev);
	u32 val;
	ssize_t ret;

	ret = kstrtouint(buf, 0, &val);
	if (ret)
		return ret;

	priv->rf_pwr_dis_val = val ? 1 : 0;
	gpio_set_value_cansleep(priv->rf_pwr_dis_gpio, priv->rf_pwr_dis_val);
	dev_info(dev, "RF power is %s now.", ((priv->rf_pwr_dis_val)==1)?"disabled":"enabled");
	return count;
}

	
static ssize_t show_rf_eeprom_wp(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_ctrl_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "RF EEPROM write protect is %s.", ((priv->rf_eeprom_wp_en_val)==0)?"disabled":"enabled");
	return sprintf(buf, "%d\n", priv->rf_eeprom_wp_en_val);
}

static ssize_t set_rf_eeprom_wp(struct device *dev,
		struct device_attribute *attr,
		const char *buf, size_t count)
{
	struct ocfema_ctrl_priv *priv = dev_get_drvdata(dev);
	u32 val;
	ssize_t ret;

	ret = kstrtouint(buf, 0, &val);
	if (ret)
		return ret;

	priv->rf_eeprom_wp_en_val = val ? 1 : 0;

	gpio_set_value_cansleep(priv->rf_eeprom_wp_en_gpio, priv->rf_eeprom_wp_en_val);
	dev_info(dev, "RF EEPROM write protect is %s now.", ((priv->rf_eeprom_wp_en_val)==0)?"disabled":"enabled");
	return count;
}

static ssize_t show_pga_pwr(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_ctrl_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "PGA power is %s.", ((priv->pga_pwr_dis_val)==0)?"disabled":"enabled");
	return sprintf(buf, "%d\n", priv->pga_pwr_dis_val);
}

static ssize_t set_pga_pwr(struct device *dev,
		struct device_attribute *attr,
		const char *buf, size_t count)
{
	struct ocfema_ctrl_priv *priv = dev_get_drvdata(dev);
	u32 val;
	ssize_t ret;

	ret = kstrtouint(buf, 0, &val);
	if (ret)
		return ret;

	priv->pga_pwr_dis_val = val ? 0 : 1;
	gpio_set_value_cansleep(priv->pga_pwr_dis_gpio, priv->pga_pwr_dis_val);
	dev_info(dev, "PGA power is %s now.", ((priv->pga_pwr_dis_val)==0)?"disabled":"enabled");
	return count;
}

static ssize_t show_pa_state(struct device *dev,
                struct device_attribute *attr, char *buf)
{
        struct ocfema_ctrl_priv *priv = dev_get_drvdata(dev);
        dev_info(dev, "PA is %s.", ((priv->pa_dis_val)==0)?"disabled":"enabled");
        return sprintf(buf, "%d\n", priv->pa_dis_val);
}

static ssize_t set_pa_state(struct device *dev,
                struct device_attribute *attr,
                const char *buf, size_t count)
{
        struct ocfema_ctrl_priv *priv = dev_get_drvdata(dev);
        u32 val;
        ssize_t ret;

        ret = kstrtouint(buf, 0, &val);
        if (ret)
                return ret;

        priv->pa_dis_val = val ? 0 : 1; //negating
        gpio_set_value_cansleep(priv->pa_dis_gpio, priv->pa_dis_val);
        dev_info(dev, "PA is in %s state now.", ((priv->pa_dis_val)==0)?"disabled":"enabled");
        return count;
}

static ssize_t show_mp_eeprom_wp(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_ctrl_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "μ-processor EEPROM write protect is %s.", ((priv->mp_eeprom_wp_en_val)==0)?"disabled":"enabled");
	return sprintf(buf, "%d\n", priv->mp_eeprom_wp_en_val);
}

static ssize_t set_mp_eeprom_wp(struct device *dev,
		struct device_attribute *attr,
		const char *buf, size_t count)
{
	struct ocfema_ctrl_priv *priv = dev_get_drvdata(dev);
	u32 val;
	ssize_t ret;

	ret = kstrtouint(buf, 0, &val);
	if (ret)
		return ret;

	priv->mp_eeprom_wp_en_val = val ? 1 : 0;

	gpio_set_value_cansleep(priv->mp_eeprom_wp_en_gpio, priv->mp_eeprom_wp_en_val);
	dev_info(dev, "μ-processor EEPROM write protect is %s now.", ((priv->mp_eeprom_wp_en_val)==0)?"disabled":"enabled");
	return count;
}



static DEVICE_ATTR(rf_pwr, S_IWUSR | S_IRUGO, show_rf_pwr, set_rf_pwr);
static DEVICE_ATTR(pga_pwr, S_IWUSR | S_IRUGO, show_pga_pwr, set_pga_pwr);
static DEVICE_ATTR(pa_state, S_IWUSR | S_IRUGO, show_pa_state, set_pa_state);
static DEVICE_ATTR(mp_eeprom_wp, S_IWUSR | S_IRUGO, show_mp_eeprom_wp, set_mp_eeprom_wp);
static DEVICE_ATTR(rf_eeprom_wp, S_IWUSR | S_IRUGO, show_rf_eeprom_wp, set_rf_eeprom_wp);

static struct attribute *ocfema_ctrl_attrs[] = {
	&dev_attr_rf_pwr.attr,
	&dev_attr_pga_pwr.attr,
	&dev_attr_pa_state.attr,
	&dev_attr_mp_eeprom_wp.attr,
	&dev_attr_rf_eeprom_wp.attr,
	NULL,
};

static const struct attribute_group ocfema_ctrl_attr_group = {
	.attrs = ocfema_ctrl_attrs,
};

int ocfema_ctrl_parse_dt(struct platform_device *pdev)
{
	int ret = 0; 
	struct device_node *np = pdev->dev.of_node;
	struct ocfema_ctrl_priv *priv = dev_get_drvdata(&pdev->dev);
	if (!np) {
		return -EINVAL;
	}
	
	/*RF Power disable */
	priv->rf_pwr_dis_gpio = of_get_named_gpio(np, "rf-pwr-dis-gpio", 0);
	if (priv->rf_pwr_dis_gpio < 0) {
		dev_err(&pdev->dev, "Can't read gpio rf-pwr-dis-gpio\n");
		return -EINVAL;
	}

	ret = of_property_read_u32(np, "rf-pwr-dis-default", &priv->rf_pwr_dis_default);
	if (ret < 0) {
		dev_err(&pdev->dev, "Can't read RF power default value from device tree.\n");
		return ret;
	}
   	priv->rf_pwr_dis_val = priv->rf_pwr_dis_default;	
	dev_info(&pdev->dev, "RF power default value is %d\n", priv->rf_pwr_dis_default);
		
	// RF board EEPROM write protect
	priv->rf_eeprom_wp_en_gpio = of_get_named_gpio(np, "rf-eeprom-wp-en-gpio", 0);
	if (priv->rf_eeprom_wp_en_gpio < 0) {
		dev_err(&pdev->dev, "Can't read gpio rf-eeprom-wp-en-gpio\n");
		return -EINVAL;
	}

	ret = of_property_read_u32(np, "rf-eeprom-wp-en-default", &priv->rf_eeprom_wp_en_default);
	if (ret < 0) {
		dev_err(&pdev->dev, "Can't read RF board EEPROM write protect default value from device tree.\n");
		return ret;
	}
	priv->rf_eeprom_wp_en_val = priv->rf_eeprom_wp_en_default;
	dev_info(&pdev->dev, "RF board EEPROM write protect default value read is %d\n", priv->rf_eeprom_wp_en_default);

	/*PGA Power disable */
	priv->pga_pwr_dis_gpio = of_get_named_gpio(np, "pga-pwr-dis-gpio", 0);
	if (priv->pga_pwr_dis_gpio < 0) {
		dev_err(&pdev->dev, "Can't read gpio pga-pwr-dis-gpio\n");
		return -EINVAL;
	}

	ret = of_property_read_u32(np, "pga-pwr-dis-default", &priv->pga_pwr_dis_default);
	if (ret < 0) {
		dev_err(&pdev->dev, "Can't read PGA power default value from device tree.\n");
		return ret;
	}
	priv->pga_pwr_dis_val = priv->pga_pwr_dis_default;
	dev_info(&pdev->dev, "PGA power default value is %d\n", priv->pga_pwr_dis_default);
	
	/*PA state disable */
        priv->pa_dis_gpio = of_get_named_gpio(np, "pa-dis-gpio", 0);
        if (priv->pa_dis_gpio < 0) {
                dev_err(&pdev->dev, "Can't read gpio pa-dis-gpio\n");
                return -EINVAL;
        }

        ret = of_property_read_u32(np, "pa-dis-default", &priv->pa_dis_default);
        if (ret < 0) {
                dev_err(&pdev->dev, "Can't read PA state default value from device tree.\n");
                return ret;
        }
        priv->pa_dis_val = priv->pa_dis_default;
        dev_info(&pdev->dev, "PA state default value is %d\n", priv->pa_dis_default);

	// μ-proceesor board EEPROM write protect
	priv->mp_eeprom_wp_en_gpio = of_get_named_gpio(np, "mp-eeprom-wp-en-gpio", 0);
	if (priv->mp_eeprom_wp_en_gpio < 0) {
		dev_err(&pdev->dev, "Can't read gpio mp-eeprom-wp-en-gpio\n");
		return -EINVAL;
	}

	ret = of_property_read_u32(np, "mp-eeprom-wp-en-default", &priv->mp_eeprom_wp_en_default);
	if (ret < 0) {
		dev_err(&pdev->dev, "Can't read μ-processor board EEPROM write protect default value from device tree.\n");
		return ret;
	}
	priv->mp_eeprom_wp_en_val = priv->mp_eeprom_wp_en_default;
	dev_info(&pdev->dev, "μ-processor board EEPROM write protect default value read is %d\n", priv->mp_eeprom_wp_en_default);

	return 0;
}

static int ocfema_ctrl_probe(struct platform_device *pdev)
{ 
	struct ocfema_ctrl_priv *priv;
	int ret;
	
    	priv = devm_kzalloc(&pdev->dev, sizeof(struct ocfema_ctrl_priv), GFP_KERNEL);
	if (!priv) {
		dev_err(priv->dev, "Unable to allocate memory.\n");
		return -ENOMEM;
	}

	priv->dev = &pdev->dev;
	platform_set_drvdata(pdev, priv);

	ret = ocfema_ctrl_parse_dt(pdev);
	if (ret) {	
		dev_err(&pdev->dev, "Parsing failed for ocfema node.");
		return ret;
	}	
	
	dev_info(&pdev->dev, "Configure default OC-FEMA controls on boot.");
   	
	/* Configuring RF power */
	sprintf(priv->rf_pwr_dis_szGpio, "rf_pwr");
	ret = gpio_request(priv->rf_pwr_dis_gpio, priv->rf_pwr_dis_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->rf_pwr_dis_gpio, priv->rf_pwr_dis_szGpio);
		return ret;
	}

	ret = gpio_direction_output(priv->rf_pwr_dis_gpio, priv->rf_pwr_dis_val);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->rf_pwr_dis_gpio, ret);
		return ret;
	}
	
	dev_info(&pdev->dev, "RF default power value %s.",(priv->rf_pwr_dis_val==1)?"disabled":"enabled");
   	
	/* Configuring RF board EEPROM */
	sprintf(priv->rf_eeprom_wp_en_szGpio, "rf_eeprom_wp");
	ret = gpio_request(priv->rf_eeprom_wp_en_gpio, priv->rf_eeprom_wp_en_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->rf_eeprom_wp_en_gpio, priv->rf_eeprom_wp_en_szGpio);
		return ret;
	}

	ret = gpio_direction_output(priv->rf_eeprom_wp_en_gpio, priv->rf_eeprom_wp_en_val);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->rf_eeprom_wp_en_gpio, ret);
		return ret;
	}

	dev_info(&pdev->dev, "RF board EEPROM write protect is %s.",(priv->rf_eeprom_wp_en_val==0)?"disabled":"enabled");	
	/* Configuring PGA power */
	sprintf(priv->pga_pwr_dis_szGpio, "pga_pwr");
	ret = gpio_request(priv->pga_pwr_dis_gpio, priv->pga_pwr_dis_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->pga_pwr_dis_gpio, priv->pga_pwr_dis_szGpio);
		return ret;
	}

	ret = gpio_direction_output(priv->pga_pwr_dis_gpio, priv->pga_pwr_dis_val);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->pga_pwr_dis_gpio, ret);
		return ret;
	}
	
	dev_info(&pdev->dev, "PGA default power value %s.",(priv->pga_pwr_dis_val==0)?"disabled":"enabled");
   	
	/* Configuring PA state */
        sprintf(priv->pa_dis_szGpio, "pa_state");
        ret = gpio_request(priv->pa_dis_gpio, priv->pa_dis_szGpio);
        if ( ret )
        {
                dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
                        priv->pa_dis_gpio, priv->pa_dis_szGpio);
                return ret;
        }

        ret = gpio_direction_output(priv->pa_dis_gpio, priv->pa_dis_val);
        if ( ret )
        {
                dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
                        priv->pa_dis_gpio, ret);
                return ret;
        }

        dev_info(&pdev->dev, "PA state default value %s.",(priv->pa_dis_val==0)?"disabled":"enabled");

	/* Configuring μ-processor board EEPROM */
	sprintf(priv->mp_eeprom_wp_en_szGpio, "mp_eeprom_wp");
	ret = gpio_request(priv->mp_eeprom_wp_en_gpio, priv->mp_eeprom_wp_en_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->mp_eeprom_wp_en_gpio, priv->mp_eeprom_wp_en_szGpio);
		return ret;
	}

	ret = gpio_direction_output(priv->mp_eeprom_wp_en_gpio, priv->mp_eeprom_wp_en_val);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->mp_eeprom_wp_en_gpio, ret);
		return ret;
	}

	dev_info(&pdev->dev, "μ-processor board EEPROM write protect is %s.",(priv->mp_eeprom_wp_en_val==0)?"disabled":"enabled");	

	dev_info(&pdev->dev, "Setting Sysfs for OC-FEMA controls.");
	
	ret = sysfs_create_group(&priv->dev->kobj, &ocfema_ctrl_attr_group);
	if (ret) {
		dev_err(priv->dev, "unable to create sysfs files\n");
		return ret;
	}
	return 0;
}

static int ocfema_ctrl_remove(struct platform_device *pdev)
{
	struct ocfema_ctrl_priv *priv = platform_get_drvdata(pdev);
	sysfs_remove_group(&priv->dev->kobj, &ocfema_ctrl_attr_group);
	return 0;
}

static const struct of_device_id ocfema_of_match[] = {
	{ .compatible = "oc,fema-ctrl", },
	{}
};
MODULE_DEVICE_TABLE(of, ocfema_of_match);

static struct platform_driver ocfema_ctrl_driver = {
	.driver = {
		.name = "ocfema_ctrl",
		.owner = THIS_MODULE,
		.of_match_table = of_match_ptr(ocfema_of_match),
	},
	.probe = ocfema_ctrl_probe,
	.remove = ocfema_ctrl_remove,
};

module_platform_driver(ocfema_ctrl_driver);

MODULE_DESCRIPTION("OCFEMA Control Driver");
MODULE_AUTHOR("<vthakur@fb.com>");
MODULE_LICENSE("GPL");
