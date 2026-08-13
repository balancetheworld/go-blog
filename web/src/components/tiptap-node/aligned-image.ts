import Image from '@tiptap/extension-image'
import { ReactNodeViewRenderer } from '@tiptap/react'
import { ResizableImageView } from '@/components/tiptap-node/resizable-image-view'

export const AlignedImage = Image.extend({
  addAttributes() {
    return {
      ...this.parent?.(),
      align: {
        default: 'center',
        parseHTML: element => element.getAttribute('data-align') ?? 'center',
        renderHTML: attributes => ({
          'data-align': attributes.align,
        }),
      },
      width: {
        default: null,
        parseHTML: (element) => {
          const value = element.getAttribute('data-width') ?? element.getAttribute('width')
          const width = Number.parseInt(value ?? '', 10)
          return Number.isFinite(width) && width > 0 ? width : null
        },
        renderHTML: attributes => attributes.width
          ? {
              'data-width': attributes.width,
              'width': attributes.width,
            }
          : {},
      },
    }
  },
  addNodeView() {
    return ReactNodeViewRenderer(ResizableImageView)
  },
  renderHTML({ HTMLAttributes }) {
    const align = HTMLAttributes['data-align'] ?? 'center'
    const style = align === 'left'
      ? 'display:block;margin-right:auto;'
      : align === 'right'
        ? 'display:block;margin-left:auto;'
        : 'display:block;margin-left:auto;margin-right:auto;'

    return ['img', {
      ...this.options.HTMLAttributes,
      ...HTMLAttributes,
      style,
    }]
  },
})
