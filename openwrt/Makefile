include $(TOPDIR)/rules.mk

PKG_NAME:=sing-box-sub
PKG_VERSION:=1.0.0
PKG_RELEASE:=1

PKG_SOURCE_PROTO:=git
PKG_SOURCE_URL:=https://github.com/yourname/sing-box-sub.git
PKG_SOURCE_VERSION:=main
PKG_MIRROR_HASH:=skip

PKG_LICENSE:=MIT
PKG_MAINTAINER:=sing-box-sub

include $(INCLUDE_DIR)/package.mk

define Package/sing-box-sub
    SECTION:=net
    CATEGORY:=Network
    TITLE:=sing-box subscription converter daemon
    DEPENDS:=+ca-bundle
    URL:=https://github.com/yourname/sing-box-sub
endef

define Package/sing-box-sub/description
    sing-box subscription converter.
    Converts clash-meta subscription URLs into sing-box configuration.
endef

define Build/Compile
    cd $(PKG_BUILD_DIR)/src && \
    GOOS=linux GOARCH=$(GO_ARCH) \
    go build -ldflags="-s -w" -o $(PKG_BUILD_DIR)/sing-box-sub .
endef

define Package/sing-box-sub/install
    $(INSTALL_DIR) $(1)/usr/bin
    $(INSTALL_BIN) $(PKG_BUILD_DIR)/sing-box-sub $(1)/usr/bin/

    $(INSTALL_DIR) $(1)/etc/init.d
    $(INSTALL_BIN) ./openwrt/sing-box-sub.init $(1)/etc/init.d/sing-box-sub

    $(INSTALL_DIR) $(1)/etc/config
    $(INSTALL_CONF) ./openwrt/sing-box-sub.config $(1)/etc/config/sing-box-sub
endef

define Package/sing-box-sub/postinst
#!/bin/sh
[ -n "$$IPKG_INSTROOT" ] && exit 0
/etc/init.d/sing-box-sub enable
/etc/init.d/sing-box-sub start
endef

$(eval $(call BuildPackage,sing-box-sub))
