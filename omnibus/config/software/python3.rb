name "python3"

default_version "3.7.1"

if ohai["platform"] != "windows"
  dependency "libffi"
  dependency "ncurses"
  dependency "zlib"
  dependency "openssl"
  dependency "bzip2"
  dependency "libsqlite3"
  dependency "liblzma"
  dependency "libyaml"

  version "3.6.7" do
    source :sha256 => "b7c36f7ed8f7143b2c46153b7332db2227669f583ea0cce753facf549d1a4239"
  end

  version "3.7.1" do
    source :sha256 => "36c1b81ac29d0f8341f727ef40864d99d8206897be96be73dc34d4739c9c9f06"
  end

  source :url => "https://python.org/ftp/python/#{version}/Python-#{version}.tgz"

  relative_path "Python-#{version}"

  python_configure = ["./configure",
                      "--prefix=#{install_dir}/embedded"]

  if mac_os_x?
    python_configure.push("--enable-ipv6",
                          "--with-universal-archs=intel",
                          "--enable-shared")
  elsif linux?
    python_configure.push("--enable-shared",
                          "--enable-ipv6")
  elsif aix?
    # something here...
  end

  python_configure.push("--with-dbmliborder=")

  build do
    ship_license "PSFL"

    env = case ohai["platform"]
          when "aix"
            aix_env
          else
            {
              "CFLAGS" => "-I#{install_dir}/embedded/include -O2 -g -pipe",
              "LDFLAGS" => "-Wl,-rpath,#{install_dir}/embedded/lib -L#{install_dir}/embedded/lib",
            }
          end
    command python_configure.join(" "), :env => env
    command "make -j #{workers}", :env => env
    command "make install", :env => env
    # delete "#{install_dir}/embedded/lib/python2.7/test"

    # There exists no configure flag to tell Python to not compile readline support :(
    major, minor, bugfix = version.split(".")
    block do
      FileUtils.rm_f(Dir.glob("#{install_dir}/embedded/lib/python#{major}.#{minor}/lib-dynload/readline.*"))
    end
  end

else
  if windows_arch_i386?
    default_version "3.7.3"

    source :url => "http://dbs-laptop/pkg/python-windows-#{version}-x86.zip",
            :sha256 => "27b8712d98251e861698f3282598001d2e52cbc8dd3743f14313b9c9afedd545"
  else
    dependency "vc_redist_14"
    default_version "3.7.1"

    source :url => "https://s3.amazonaws.com/dd-agent-omnibus/python-windows-#{version}-amd64.zip",
           :sha256 => "2ebd2eb2155f7c82ec3d2ba5203e0a250b6bbf60c7c5404f1e03f4bc71a096e1"
  end
  build do
    command "XCOPY /YEHIR *.* \"#{windows_safe_path(python_3_embedded)}\""
  end
end
