import React, { useState } from 'react';
import { saveAs } from 'file-saver';
import './JsonUploader.css';

const JsonUploader = () => {
  const [jsonData, setJsonData] = useState(null);
  const [flatKeys, setFlatKeys] = useState([]);
  const [jsonFiles, setJsonFiles] = useState({});
  const [selectedKeys, setSelectedKeys] = useState({});
  const [fileNames, setFileNames] = useState({});

  const handleFileUpload = (event, key) => {
    const file = event.target.files[0];
    const reader = new FileReader();

    reader.onload = (e) => {
      const json = JSON.parse(e.target.result);
      setJsonFiles((prevJsonFiles) => ({
        ...prevJsonFiles,
        [key]: json,
      }));
      const keys = collectFlatKeys(json);
      setSelectedKeys((prevSelectedKeys) => ({
        ...prevSelectedKeys,
        [key]: keys[0], // Выбираем первый ключ по умолчанию
      }));
      setFileNames((prevFileNames) => ({
        ...prevFileNames,
        [key]: file.name,
      }));
    };

    reader.readAsText(file);
  };

  const collectFlatKeys = (data, prefix = '') => {
    let keys = [];
    if (typeof data === 'object' && data !== null) {
      for (let key in data) {
        if (typeof data[key] === 'object' && data[key] !== null) {
          keys = keys.concat(collectFlatKeys(data[key], `${prefix}${key}.`));
        } else {
          keys.push(`${prefix}${key}`);
        }
      }
    }
    return keys;
  };

  const handleKeyChange = (key, value) => {
    setSelectedKeys((prevSelectedKeys) => ({
      ...prevSelectedKeys,
      [key]: value,
    }));
  };

  const handleSecondFileUpload = (event) => {
    const file = event.target.files[0];
    const reader = new FileReader();

    reader.onload = (e) => {
      const json = JSON.parse(e.target.result);
      setJsonData(json);
      const keys = collectFlatKeys(json);
      setFlatKeys(keys);
      setJsonFiles(keys.reduce((acc, key) => {
        acc[key] = null;
        return acc;
      }, {}));
    };

    reader.readAsText(file);
  };

  const generateNewJson = () => {
    const newJson = {};
    flatKeys.forEach((key) => {
      if (selectedKeys[key]) {
        newJson[key] = {
          key: selectedKeys[key],
          fileName: fileNames[key],
        };
      }
    });

    const blob = new Blob([JSON.stringify(newJson, null, 2)], { type: 'application/json' });
    saveAs(blob, 'new_json.json');
  };

  return (
    <div className="json-uploader">
      <input type="file" accept=".json" onChange={handleSecondFileUpload} />
      {jsonData && flatKeys.length > 0 && (
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>Требуемые данные</th>
                <th>Выбор БД</th>
                <th>Поле</th>
                <th>Параметры AI</th>
              </tr>
            </thead>
            <tbody>
              {flatKeys.map((key, index) => (
                <tr key={index}>
                  <td className="required-data">{key}</td>
                  <td>
                    <input type="file" accept=".json" onChange={(e) => handleFileUpload(e, key)} />
                  </td>
                  <td>
                    {jsonFiles[key] && (
                      <select
                        value={selectedKeys[key]}
                        onChange={(e) => handleKeyChange(key, e.target.value)}
                      >
                        <option value="">Выберите ключ</option>
                        {collectFlatKeys(jsonFiles[key]).map((firstKey, firstIndex) => (
                          <option key={firstIndex} value={firstKey}>
                            {firstKey}
                          </option>
                        ))}
                      </select>
                    )}
                  </td>
                  <td>
                    {/* Параметры AI можно добавить позже */}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <button onClick={generateNewJson} className="generate-button">Сформировать новый JSON</button>
        </div>
      )}
    </div>
  );
};

export default JsonUploader;

