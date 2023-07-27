# Makefile:

# Generate Resources
bundle_images:
	 fyne bundle --package commons -o ./commons/images.go ./resources

package_mac_gui:
	fyne package -os darwin --name hubPower
