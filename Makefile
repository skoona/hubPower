# Makefile:

# Generate Resources
bundle_images:
	 fyne bundle --package commons -o ./internal/commons/images.go ./internal/commons/resources

package_mac_gui:
	fyne package -os darwin --name hubPower
