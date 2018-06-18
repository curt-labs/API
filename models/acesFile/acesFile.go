package acesFile

import (
	"bytes"
	"os"
	"time"

	"github.com/curt-labs/API/models/brand"
	"github.com/pkg/errors"

	"github.com/jlaffaye/ftp"
)

type FtpConfig struct {
	Address    string
	User       string
	Password   string
	Connection *ftp.ServerConn
}

var (
	ftpAddr = os.Getenv("FTP_HOST")
	ftpUser = os.Getenv("FTP_USERNAME")
	ftpPass = os.Getenv("FTP_PASSWORD")
)

func GetAcesFile(brand brand.Brand, version string) (string, error) {
	//Establish connection to the FTP
	ftpConfig := FtpConfig{
		Address:  ftpAddr,
		User:     ftpUser,
		Password: ftpPass,
	}

	err := ftpConfig.NewConnection()
	if err != nil {
		return "", err
	}
	defer ftpConfig.Connection.Quit()

	var fileName string
	path := "/Vendor Login Files/masterlibrary/01resources/ACES/" + version

	//set filename depending on brand
	if brand.ID == 1 {
		path += "/CURT"
		fileName = "CURT_ACES" + version + ".xml"
	} else if brand.ID == 3 {
		path += "/ARIES"
		fileName = "ARIES_ACES" + version + ".xml"
	} else if brand.ID == 4 {
		path += "/LUVERNE"
		fileName = "LUVERNE_ACES" + version + ".xml"
	} else if brand.ID == 6 {
		path += "/UWS"
		fileName = "UWS_ACES" + version + ".xml"
	}

	file, err := ftpConfig.GetFile(path, fileName)
	if err != nil {
		return "", err
	}

	return file, nil
}

// NewConnection establishes a new connection the provided FTP server.
func (f *FtpConfig) NewConnection() error {
	conn, err := ftp.DialTimeout(f.Address, time.Second*2)
	if err != nil {
		return err
	}

	err = conn.Login(f.User, f.Password)
	if err != nil {
		conn.Quit()
		return err
	}

	f.Connection = conn

	return nil
}

// GetFile executes a RETR command on the FTP, pulls down the file and
// parses it as an XML file.
func (f *FtpConfig) GetFile(path, fileName string) (string, error) {
	if f.Connection == nil {
		return "", errors.New("no connection established")
	}

	err := f.Connection.ChangeDir(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to change directory %s", path)
	}

	cl, err := f.Connection.Retr(fileName)
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve file %s", fileName)
	}
	defer cl.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(cl)
	file := buf.String()

	return file, nil
}
