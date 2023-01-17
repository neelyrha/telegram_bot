package files

import (
	"GoRestApi/lib/e"
	"GoRestApi/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }() //способ обработки ошибок

	fPath := filepath.Join(s.basePath, page.UserName) //путь в который будет сохраняться файл

	if err := os.Mkdir(fPath, defaultPerm); err != nil {
		return err
	} //создаем нужные директории в этом пути

	fName, err := fileName(page) //формируем имя файла
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName) //дописываем имя файла к пути

	file, err := os.Create(fPath) //создаем файл
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }() //для того что бы обработать ошибку которую мы не хотим обрабатывать

	if err := gob.NewEncoder(file).Encode(page); err != nil { //записываем в файл страничку в нужном формате
		return err
	}
	return nil
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)

	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprint("can't remove file: %s", path)
		return e.Wrap(msg, err)
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}
	rand.Seed(time.Now().UnixNano()) //генерируем случайное число засчет сида-времени. что так
	n := rand.Intn(len(files))       //файл это и есть сохраненная страничка

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))

}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)

	if err != nil {
		return false, e.Wrap("can't find file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)
		return false, e.Wrap(msg, err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
