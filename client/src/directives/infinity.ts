export function infinity(node: Node) {
  let wrapper: HTMLElement | null;
  let content: HTMLElement | null;
  let block: HTMLElement | null;

  // node.dispatchEvent(
  //   new CustomEvent("panstart", {
  //     detail: { x, y },
  //   })
  // );

  const x: HTMLElement | null = document.getElementById("wrapper");
  const y: HTMLElement | null = document.getElementById("content");
  const z: HTMLElement | null = document.getElementById("block");

  console.log(node);
  wrapper = x;
  content = y;
  block = z;

  const handleScroll = (): void => {
    if (wrapper!.scrollTop + wrapper!.offsetHeight > content!.offsetHeight) {
      wrapper!.scrollTop = block!.offsetHeight - wrapper!.offsetHeight;
    }
  };

  node.addEventListener("scroll", handleScroll);

  return {
    destroy() {
      node.removeEventListener("scroll", handleScroll);
    },
  };
}
