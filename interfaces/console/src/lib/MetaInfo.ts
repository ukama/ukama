const getMetaInfo = async () => {
  return await fetch('', {
    method: 'GET',
  })
    .then((response) => response.text())
    .then((data) => JSON.parse(data))
    .then((data) =>
      fetch(`/${data.ip}/json/`, {
        method: 'GET',
      }),
    )
    .then((response) => response.text())
    .then((data) => JSON.parse(data))
    .catch((err) => {
      console.log(err);
      return {};
    });
};
