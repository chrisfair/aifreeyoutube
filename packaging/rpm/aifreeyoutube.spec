# ==============================================================================
# AI-Free YouTube: RPM Package Specification
# ==============================================================================
# Target distributions: Fedora 38+, Red Hat Enterprise Linux 9+, openSUSE Leap/Tumbleweed
# Build command:
#   rpmbuild -bb packaging/rpm/aifreeyoutube.spec --define "_topdir $(pwd)/build/rpm"
# ==============================================================================

Name:           aifreeyoutube
Version:        1.0.0
Release:        1%{?dist}
Summary:        Authentic human video portal filtering out AI slop on YouTube

License:        MIT
URL:            https://github.com/christopherpatrickfair/aifreeyoutube
Source0:        %{name}-%{version}.tar.gz

BuildArch:      x86_64
Requires:       webkit2gtk4.1 >= 2.38.0
Requires:       gtk3
Requires:       hicolor-icon-theme
Requires:       glib2

%description
AI-Free YouTube transparently intercepts and filters YouTube search and trending
content using an explainable, multi-signal heuristic engine. It quarantines
automated ChatGPT scripts, synthetic voiceover content mills, and clickbait
alarmism formulas, delivering an authentic human viewing experience.
Supports both native WebKit2GTK desktop windows and local web server mode.

%prep
# No prep needed for precompiled binary packages

%build
# Binaries are compiled via Wails and Go standard tooling prior to rpmbuild

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT%{_bindir}
mkdir -p $RPM_BUILD_ROOT%{_datadir}/applications
mkdir -p $RPM_BUILD_ROOT%{_datadir}/icons/hicolor/scalable/apps
mkdir -p $RPM_BUILD_ROOT%{_datadir}/icons/hicolor/512x512/apps
mkdir -p $RPM_BUILD_ROOT%{_unitdir}

# Install executables
install -m 755 %{_sourcedir}/aifreeyoutube $RPM_BUILD_ROOT%{_bindir}/aifreeyoutube
install -m 755 %{_sourcedir}/aifreeyoutube-server $RPM_BUILD_ROOT%{_bindir}/aifreeyoutube-server

# Install desktop entry and icons
install -m 644 %{_sourcedir}/aifreeyoutube.desktop $RPM_BUILD_ROOT%{_datadir}/applications/aifreeyoutube.desktop
install -m 644 %{_sourcedir}/appicon.svg $RPM_BUILD_ROOT%{_datadir}/icons/hicolor/scalable/apps/aifreeyoutube.svg
install -m 644 %{_sourcedir}/appicon.png $RPM_BUILD_ROOT%{_datadir}/icons/hicolor/512x512/apps/aifreeyoutube.png

# Install systemd service
install -m 644 %{_sourcedir}/aifreeyoutube.service $RPM_BUILD_ROOT%{_unitdir}/aifreeyoutube.service

%post
/bin/touch --no-create %{_datadir}/icons/hicolor &>/dev/null || :
if [ $1 -eq 1 ] ; then
    # Initial installation
    /usr/bin/update-desktop-database &> /dev/null || :
fi

%postun
if [ $1 -eq 0 ] ; then
    # Complete removal
    /usr/bin/update-desktop-database &> /dev/null || :
    /bin/touch --no-create %{_datadir}/icons/hicolor &>/dev/null || :
    %{_bindir}/gtk-update-icon-cache %{_datadir}/icons/hicolor &>/dev/null || :
fi

%posttrans
%{_bindir}/gtk-update-icon-cache %{_datadir}/icons/hicolor &>/dev/null || :

%files
%{_bindir}/aifreeyoutube
%{_bindir}/aifreeyoutube-server
%{_datadir}/applications/aifreeyoutube.desktop
%{_datadir}/icons/hicolor/scalable/apps/aifreeyoutube.svg
%{_datadir}/icons/hicolor/512x512/apps/aifreeyoutube.png
%{_unitdir}/aifreeyoutube.service

%changelog
* Mon Oct 12 2026 Christopher Patrick Fair <christopherpatrickfair@gmail.com> - 1.0.0-1
- Initial release of AI-Free YouTube with Wails v2 desktop integration and heuristic filter engine.
