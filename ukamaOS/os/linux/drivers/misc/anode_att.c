 #include <linux/err.h>
 #include <linux/module.h>
 #include <linux/types.h>
 #include <linux/platform_device.h>
 #include <linux/io.h>
 #include <linux/of_gpio.h>
 #include <linux/of.h>
 #include <linux/sysfs.h>
 #include <linux/delay.h>

#define NR_ATT_GPIO                     6

struct ocfema_att_priv
{
    	struct device *dev;
	
	int tx_att_le_gpio;
	char tx_att_le_szGpio[32];
	u8 tx_att_le_val;
	u32 tx_att_le_default;
	
	int tx_att_gpio[NR_ATT_GPIO];
	char tx_att_szGpio[NR_ATT_GPIO][32];
	u8 tx_att_val;
	u32 tx_att_default;
	
	int rx_att_le_gpio;
	char rx_att_le_szGpio[32];
	u8 rx_att_le_val;
	u32 rx_att_le_default;
	
	int rx_att_gpio[NR_ATT_GPIO];
	char rx_att_szGpio[NR_ATT_GPIO][32];
	u8 rx_att_val;
	u32 rx_att_default;
	
};

static ssize_t show_rx_att_le(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_att_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "RX attenuation latch is %s.", ((priv->rx_att_le_val)==0)?"disabled":"enabled");
	return sprintf(buf, "%d\n", priv->rx_att_le_val);
}

static ssize_t set_rx_att_le(struct device *dev,
		struct device_attribute *attr,
		const char *buf, size_t count)
{
	struct ocfema_att_priv *priv = dev_get_drvdata(dev);
	u32 val;
	ssize_t ret;

	ret = kstrtouint(buf, 0, &val);
	if (ret)
		return ret;

	priv->rx_att_le_val = val ? 1 : 0;
	gpio_set_value_cansleep(priv->rx_att_le_gpio, priv->rx_att_le_val);
	dev_info(dev, "RX attenuation latch is %s now.", ((priv->rx_att_le_val)==0)?"disabled":"enabled");
	return count;
}

static ssize_t show_rx_att(struct device *dev,
                struct device_attribute *attr, char *buf)
{
        struct ocfema_att_priv *priv = dev_get_drvdata(dev);
        return sprintf(buf, "%d\n", priv->rx_att_val);
}

static ssize_t set_rx_att(struct device *dev,
                struct device_attribute *attr,
                const char *buf, size_t count)
{
        struct ocfema_att_priv *priv = dev_get_drvdata(dev);
        u32 att;
        ssize_t ret;
        int i;

        ret = kstrtouint(buf, 0, &att);
        if (ret)
                return ret;

        priv->rx_att_val = (att < 64) ? att : 63;

        for (i = 0; i < NR_ATT_GPIO; i++)
        {
                gpio_set_value_cansleep(priv->rx_att_gpio[i], (priv->rx_att_val & (1 << i)));
        }
        return count;
}

	
static ssize_t show_tx_att_le(struct device *dev,
		struct device_attribute *attr, char *buf)
{
	struct ocfema_att_priv *priv = dev_get_drvdata(dev);
	dev_info(dev, "RX attenuation latch is %s.", ((priv->rx_att_le_val)==0)?"disabled":"enabled");
	return sprintf(buf, "%d\n", priv->tx_att_le_val);
}

static ssize_t set_tx_att_le(struct device *dev,
		struct device_attribute *attr,
		const char *buf, size_t count)
{
	struct ocfema_att_priv *priv = dev_get_drvdata(dev);
	u32 val;
	ssize_t ret;

	ret = kstrtouint(buf, 0, &val);
	if (ret)
		return ret;

	priv->tx_att_le_val = val ? 1 : 0;

	gpio_set_value_cansleep(priv->tx_att_le_gpio, priv->tx_att_le_val);
	dev_info(dev, "RX attenuation latch is %s now.", ((priv->rx_att_le_val)==0)?"disabled":"enabled");
	return count;
}

static ssize_t show_tx_att(struct device *dev,
                struct device_attribute *attr, char *buf)
{
        struct ocfema_att_priv *priv = dev_get_drvdata(dev);
        return sprintf(buf, "%d\n", priv->tx_att_val);
}

static ssize_t set_tx_att(struct device *dev,
                struct device_attribute *attr,
                const char *buf, size_t count)
{
        struct ocfema_att_priv *priv = dev_get_drvdata(dev);
        u32 att;
        ssize_t ret;
        int i;

        ret = kstrtouint(buf, 0, &att);
        if (ret)
                return ret;

        priv->tx_att_val = (att < 64) ? att : 63;

        for (i = 0; i < NR_ATT_GPIO; i++)
        {
                gpio_set_value_cansleep(priv->tx_att_gpio[i], (priv->tx_att_val & (1 << i)));
        }
        return count;
}

static DEVICE_ATTR(rx_att_le, S_IWUSR | S_IRUGO, show_rx_att_le, set_rx_att_le);
static DEVICE_ATTR(rx_att, S_IWUSR | S_IRUGO, show_rx_att, set_rx_att);

static DEVICE_ATTR(tx_att_le, S_IWUSR | S_IRUGO, show_tx_att_le, set_tx_att_le);
static DEVICE_ATTR(tx_att, S_IWUSR | S_IRUGO, show_tx_att, set_tx_att);

static struct attribute *ocfema_attrs[] = {
	&dev_attr_rx_att_le.attr,
	&dev_attr_rx_att.attr,
	&dev_attr_tx_att_le.attr,
	&dev_attr_tx_att.attr,
	NULL,
};

static const struct attribute_group ocfema_attr_group = {
	.attrs = ocfema_attrs,
};

int ocfema_att_parse_dt(struct platform_device *pdev)
{
	int ret = 0; 
	int nr = 0;
	int i = 0; 
	struct device_node *np = pdev->dev.of_node;
	struct ocfema_att_priv *priv = dev_get_drvdata(&pdev->dev);
	if (!np) {
		return -EINVAL;
	}
	
	/* RX attenuator latch enable */
	priv->rx_att_le_gpio = of_get_named_gpio(np, "rx-att-le-gpio", 0);
	if (priv->rx_att_le_gpio < 0) {
		dev_err(&pdev->dev, "Can't read gpio rx-att-le-gpio\n");
		return -EINVAL;
	}

	/* RX attenuator latch enable default value */
	ret = of_property_read_u32(np, "rx-att-le-default", &priv->rx_att_le_default);
	if (ret < 0) {
		dev_err(&pdev->dev, "Can't read RX att latch enable default value from DT.\n");
		return ret;
	}
	
	dev_info(&pdev->dev, "RX attenuation LE value read is %d\n", priv->rx_att_le_default);
	priv->rx_att_le_val = priv->rx_att_le_default;
	
	/* RX Attenuation setting */
	nr = of_gpio_named_count(np, "rx-att-gpios");
	if ( nr != NR_ATT_GPIO ) {
		dev_err(&pdev->dev, "Can't read all gpio required for rx-att-gpio from device tree\n"); 
		return -EINVAL;
	}
	for (i = 0; i < NR_ATT_GPIO; i++) {
		priv->rx_att_gpio[i] = of_get_named_gpio(np, "rx-att-gpios", i);
		if (priv->rx_att_gpio[i] < 0) {
			dev_err(&pdev->dev, "Can't get gpio rx-att-gpios%d\n", i);
			return -EINVAL;
		}
	}

	/* RX default attenuation */
	ret = of_property_read_u32(np, "rx-att-default", &priv->rx_att_default);
	if (ret < 0) {
		dev_err(&pdev->dev, "Can't read RX default attenuation from device tree.\n");
		return ret;
	}
	priv->rx_att_val = priv->rx_att_default;
	dev_info(&pdev->dev, "RX attenuation value read is %d\n", priv->rx_att_default);	
	
	/* TX attenuator latch enable */
	priv->tx_att_le_gpio = of_get_named_gpio(np, "tx-att-le-gpio", 0);
	if (priv->tx_att_le_gpio < 0) {
		dev_err(&pdev->dev, "Can't get gpio tx-att-le-gpio\n");
		return -EINVAL;
	}
	
	/* TX attenuator latch enable default value*/
	ret = of_property_read_u32(np, "tx-att-le-default", &priv->tx_att_le_default);
	if (ret < 0) {
		dev_err(&pdev->dev, "Can't read TX att latch enable default value from DT.\n");
		return ret;
	}
	priv->tx_att_le_val = priv->tx_att_le_default;
	dev_info(&pdev->dev, "TX zttenuation LE value read is %d\n", priv->tx_att_le_default);
	
	/* TX Attenuation setting */
	nr = of_gpio_named_count(np, "tx-att-gpios");
	if ( nr != NR_ATT_GPIO ) {
		dev_err(&pdev->dev, "Can't read all gpio required for tx-att-gpio from device tree\n");
		return -EINVAL;
	}
	for (i = 0; i < NR_ATT_GPIO; i++) {
		priv->tx_att_gpio[i] = of_get_named_gpio(np, "tx-att-gpios", i);
		if (priv->tx_att_gpio[i] < 0) {
			dev_err(&pdev->dev, "Can't get gpio tx-att-gpios%d\n", i);
			return -EINVAL;
		}
	}

	// TX PA default attenuation
	ret = of_property_read_u32(np, "tx-att-default", &priv->tx_att_default);
	if (ret < 0) {
		dev_err(&pdev->dev, "Can't read TX default attenuation from device tree.\n");
		return ret;
	}
	priv->tx_att_val = priv->tx_att_default;
	dev_info(&pdev->dev, "TX attenuation value read is %d\n", priv->tx_att_default);
	return 0;
}

static int ocfema_att_probe(struct platform_device *pdev)
{ 
	struct ocfema_att_priv *priv;
	int ret;
	int i = 0;	
    	priv = devm_kzalloc(&pdev->dev, sizeof(struct ocfema_att_priv), GFP_KERNEL);
	if (!priv) {
		dev_err(priv->dev, "Unable to allocate memory.\n");
		return -ENOMEM;
	}

	priv->dev = &pdev->dev;
	platform_set_drvdata(pdev, priv);

	ret = ocfema_att_parse_dt(pdev);
	if (ret) {	
		dev_err(&pdev->dev, "Parsing failed for ocfema-atten node.");
		return ret;
	}
	
   	/* RX attenuator latch enable */

	sprintf(priv->rx_att_le_szGpio, "rx_pa_att_le");
	ret = gpio_request(priv->rx_att_le_gpio, priv->rx_att_le_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->rx_att_le_gpio, priv->rx_att_le_szGpio);
		return ret;
	}

	ret = gpio_direction_output(priv->rx_att_le_gpio, priv->rx_att_le_val);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->rx_att_le_gpio, ret);
		return ret;
	}
	dev_info(&pdev->dev, "Setting default RX attenuation latch enable value to %d.",priv->rx_att_le_val);
	
	/* RX attenuation config */
	for (i = 0; i < NR_ATT_GPIO; i++)
        {
                // Get GPIO
                sprintf(priv->rx_att_szGpio[i], "rx_att%d", i);
                ret = gpio_request(priv->rx_att_gpio[i], priv->rx_att_szGpio[i]);
                if ( ret )
                {
                        dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
                                priv->rx_att_gpio[i], priv->rx_att_szGpio[i]);
                        return ret;
                }

                // Set direction and default output value
                ret = gpio_direction_output(priv->rx_att_gpio[i], (priv->rx_att_val & (1 << i)));
                if ( ret )
                {
                        dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
                                priv->rx_att_gpio[i], ret);
                        return ret;
                }
	}	
	dev_info(&pdev->dev, "Setting default RX attenuation value to %d.",priv->rx_att_val);

	/* TX attenuator latch enable */
	sprintf(priv->tx_att_le_szGpio, "tx_att_le");
	ret = gpio_request(priv->tx_att_le_gpio, priv->tx_att_le_szGpio);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
			priv->tx_att_le_gpio, priv->tx_att_le_szGpio);
		return ret;
	}

	ret = gpio_direction_output(priv->tx_att_le_gpio, priv->tx_att_le_val);
	if ( ret )
	{
		dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
			priv->tx_att_le_gpio, ret);
		return ret;
	}
	dev_info(&pdev->dev, "Setting default TX attenuation latch enable value to %d.",priv->tx_att_le_val);	
   		
	/* TX attenuation config */
	for (i = 0; i < NR_ATT_GPIO; i++)
        {
                // Get GPIO
                sprintf(priv->tx_att_szGpio[i], "tx_att%d", i);
                ret = gpio_request(priv->tx_att_gpio[i], priv->tx_att_szGpio[i]);
                if ( ret )
                {
                        dev_err(&pdev->dev, "Could not obtain GPIO %d: %s\n",
                                priv->tx_att_gpio[i], priv->tx_att_szGpio[i]);
                         return ret;
                }

                // Set direction and default output value
                ret = gpio_direction_output(priv->tx_att_gpio[i], (priv->tx_att_val & (1 << i)));
                if ( ret )
                {
                        dev_err(&pdev->dev, "Could not configure GPIO %d direction: %d\n",
                                priv->tx_att_gpio[i], ret);
                        return ret;
                }
        }
	dev_info(&pdev->dev, "Setting TX attenuation value to %d.",priv->tx_att_val);
	
	/* Using latch enable programming
	   Hold low while configuration and than give a 
	   high to low transistion for enabling after 10 nano secs .*/
	 
	gpio_set_value_cansleep(priv->tx_att_le_gpio, 1);
	gpio_set_value_cansleep(priv->rx_att_le_gpio, 1);
	udelay(1);
	gpio_set_value_cansleep(priv->tx_att_le_gpio, 0);
	gpio_set_value_cansleep(priv->tx_att_le_gpio, 0);	
	dev_info(&pdev->dev, "Setting Sysfs for OCFEMA Attenuations.");

	ret = sysfs_create_group(&priv->dev->kobj, &ocfema_attr_group);
	if (ret) {
		dev_err(priv->dev, "unable to create sysfs files\n");
		return ret;
	}
	return 0;
}

static int ocfema_att_remove(struct platform_device *pdev)
{
	struct ocfema_att_priv *priv = platform_get_drvdata(pdev);

	sysfs_remove_group(&priv->dev->kobj, &ocfema_attr_group);
	return 0;
}

static const struct of_device_id ocfema_of_match[] = {
	{ .compatible = "oc,fema-att", },
	{}
};
MODULE_DEVICE_TABLE(of, ocfema_of_match);

static struct platform_driver ocfema_att_driver = {
	.driver = {
		.name = "ocfema_att",
		.owner = THIS_MODULE,
		.of_match_table = of_match_ptr(ocfema_of_match),
	},
	.probe = ocfema_att_probe,
	.remove = ocfema_att_remove,
};

module_platform_driver(ocfema_att_driver);

MODULE_DESCRIPTION("OCFEMA Attenuation Driver");
MODULE_AUTHOR("<vthakur@fb.com>");
MODULE_LICENSE("GPL");

		
